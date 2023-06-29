package proxy

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/dxhbiz/go-ntrip-proxy/pkg/config"
	"github.com/dxhbiz/go-ntrip-proxy/pkg/kit/log"
)

var (
	NtripOK = []byte("200 OK")
)

const (
	ConnectTimeout = 5 * time.Second
	ReconnectTime  = 1 * time.Second
)

type Signal int

const (
	SignalLength         = 32
	ConnectSignal Signal = iota
	ConnectedSignal
	CloseSignal
	ClosedSignal
	ReconnectSignal
)

type Puller struct {
	name        string
	host        string
	port        uint16
	username    string
	password    string
	mountpoint  string
	conn        net.Conn
	isReady     bool
	isClosed    bool
	readBuffer  chan []byte
	writeBuffer chan []byte
	singalChan  chan Signal

	sendTime    int64
	receiveTime int64
	status      int
}

func newPuller(cfg config.CasterConfig) *Puller {
	p := &Puller{
		name:        cfg.Name,
		host:        cfg.Host,
		port:        cfg.Port,
		username:    cfg.Username,
		password:    cfg.Password,
		mountpoint:  cfg.Mountpoint,
		readBuffer:  make(chan []byte, ReadChanLength),
		writeBuffer: make(chan []byte, WirteChanLength),
		singalChan:  make(chan Signal, SignalLength),
		isClosed:    true,
		isReady:     false,
	}

	go p.loop()

	return p
}

func (p *Puller) connect() {
	address := fmt.Sprintf("%s:%d", p.host, p.port)
	var err error
	p.conn, err = net.DialTimeout("tcp", address, ConnectTimeout)
	if err != nil {
		log.Warnf("%s create connection error %s", p.name, err.Error())
		p.singalChan <- ReconnectSignal
		return
	}
	log.Infof("%s connection connected", p.name)

	p.singalChan <- ConnectedSignal

	go p.read()
}

func (p *Puller) read() {
	for {
		buf := make([]byte, ReadDefaultLength)
		n, err := p.conn.Read(buf)
		if err != nil {
			log.Warnf("%s puller read data error: %s", p.name, err.Error())
			p.singalChan <- ReconnectSignal
			return
		}

		if n > 0 {
			p.readBuffer <- buf[0:n]
		}
	}
}

func (p *Puller) reconnect() {
	time.Sleep(ReconnectTime)
	p.connect()
}

func (p *Puller) loop() {
	defer func() {
	}()

	for {
		select {
		case signal := <-p.singalChan:
			p.doSignal(signal)
		case rBuf := <-p.readBuffer:
			p.parseData(rBuf)
		case wBuf := <-p.writeBuffer:
			p.sendData(wBuf)
		}
	}
}

func (p *Puller) doSignal(signal Signal) {
	if signal == ConnectSignal {
		p.status += 1
		log.Infof("%s puller status %d", p.name, p.status)
		if p.isClosed {
			p.connect()
		}
	}

	if signal == ConnectedSignal {
		p.isClosed = false
		p.doRequest()
	}

	if signal == ReconnectSignal {
		if p.isClosed {
			return
		}

		p.reconnect()
	}

	if signal == CloseSignal {
		p.status -= 1
		log.Infof("%s puller status %d", p.name, p.status)
	}
}

func (p *Puller) doRequest() {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("GET /%s HTTP/1.1\r\n", p.mountpoint))
	sb.WriteString("User-Agent: Ntrip NtripSDK/1.0.0\r\n")

	userpass := fmt.Sprintf("%s:%s", p.username, p.password)
	encodeStr := base64.StdEncoding.EncodeToString([]byte(userpass))
	sb.WriteString(fmt.Sprintf("Authorization: Basic %s\r\n\r\n", encodeStr))

	data := sb.String()
	sb.Reset()
	p.sendData([]byte(data))
}

func (p *Puller) parseData(data []byte) {
	if p.isReady {
		p.receiveTime = time.Now().Unix()
		if p.receiveTime-p.sendTime >= 60 {
			log.Warnf("%s puller no nmea data has been sent for more than 60 seconds", p.name)
			if p.status == 0 {
				p.isClosed = true
				p.conn.Close()
			}
			return
		}
		publish(p.name, data)
		return
	}

	idx := bytes.Index(data, NtripOK)
	if idx >= 0 {
		p.isReady = true
		log.Infof("%s puller is ready", p.name)
	}
}

func (p *Puller) sendData(buf []byte) {
	if p.isClosed {
		return
	}
	p.sendTime = time.Now().Unix()

	_, err := p.conn.Write(buf)
	if err != nil {
		log.Errorf("%s puller write data error: %s", p.name, err.Error())
		return
	}
}

func (p *Puller) send(data []byte) {
	log.Infof("%s puller send nmea %s", p.name, string(data))
	p.writeBuffer <- data
}

func (p *Puller) changeSignal(signal Signal) {
	p.singalChan <- signal
}
