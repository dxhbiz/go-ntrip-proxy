package proxy

import (
	"bytes"
	"net"
	"strings"
	"sync"

	"github.com/dxhbiz/go-ntrip-proxy/pkg/kit/log"
	"github.com/google/uuid"
)

var (
	Prefix         = []byte("GET")
	SplitDelimiter = "\r\n"
	EndDelimiter   = []byte("\r\n\r\n")
	IcyOk          = []byte("ICY 200 OK\r\n")
	Master         = "Master"
	Slave          = "Slave"
)

const (
	MAX_HEADER_LENGTH = 20
)

type ntripClient struct {
	conn        net.Conn
	id          string
	isClosed    bool
	closeChan   chan struct{}
	readBuffer  chan []byte
	writeBuffer chan []byte
	once        sync.Once
	isReady     bool
	buf         []byte
	isMain      bool
	mountpoint  string
}

func newNtripClient(conn net.Conn) *ntripClient {
	nc := &ntripClient{
		conn:        conn,
		closeChan:   make(chan struct{}),
		readBuffer:  make(chan []byte, ReadChanLength),
		writeBuffer: make(chan []byte, WirteChanLength),
	}

	nc.id = uuid.NewString()
	log.Infof("%s new ntrip client", nc.id)

	go nc.read()

	return nc
}

func (nc *ntripClient) read() {
	for {
		if nc.isClosed {
			return
		}

		buf := make([]byte, ReadDefaultLength)
		n, err := nc.conn.Read(buf)
		if err != nil {
			log.Warnf("%s read data error: %s", nc.id, err)
			nc.callClose()
			return
		}

		if n > 0 {
			nc.readBuffer <- buf[0:n]
		}
	}
}

func (nc *ntripClient) doClose() {
	if nc.isReady {
		unSubscribe(nc.mountpoint, nc.id, nc.isMain)
		nc.isReady = false
	}
}

func (nc *ntripClient) run() {
	defer func() {
		// todo close event
		close(nc.closeChan)
		close(nc.readBuffer)
		close(nc.writeBuffer)

		nc.doClose()
	}()

	for {
		select {
		case rBuf := <-nc.readBuffer:
			nc.parseData(rBuf)
		case wBuf := <-nc.writeBuffer:
			nc.sendData(wBuf)
		case <-nc.closeChan:
			nc.isClosed = true
			log.Infof("%s got close channel", nc.id)
			return
		}
	}
}

func (nc *ntripClient) parseData(buf []byte) {
	nc.buf = append(nc.buf, buf...)
	if nc.isReady {
		nc.parseBody()
		return
	}

	nc.parseHeader()
}

func (nc *ntripClient) parseHeader() {
	if len(nc.buf) < len(Prefix) {
		return
	}

	pidx := bytes.Index(nc.buf, Prefix)
	if pidx != 0 {
		log.Warnf("%s bad request", nc.id)
		nc.conn.Close()
		return
	}

	hidx := bytes.Index(nc.buf, EndDelimiter)
	if hidx == -1 {
		if len(nc.buf) >= MAX_HEADER_LENGTH {
			log.Warnf("%s header data out of limit", nc.id)
			nc.conn.Close()
		}
		return
	}

	headerBuf := nc.buf[:hidx]
	hidx += len(EndDelimiter)
	nc.buf = nc.buf[hidx:]

	headerStr := string(headerBuf)
	headerArr := strings.Split(headerStr, SplitDelimiter)

	if len(headerArr) == 0 {
		log.Warnf("%s can't split header data", nc.id)
		nc.conn.Close()
		return
	}

	mountReqArr := strings.Split(headerArr[0], " ")
	if len(mountReqArr) != 3 {
		log.Warnf("%s can't split mountpoint", nc.id)
		nc.conn.Close()
		return
	}

	mountpoint := strings.Replace(mountReqArr[1], "/", "", 1)
	mountArr := strings.Split(mountpoint, "-")
	if len(mountArr) != 2 {
		log.Warnf("%s invalid mountpoint, it must be contain '-'", nc.id)
		nc.conn.Close()
		return
	}

	if mountArr[1] == Master {
		nc.isMain = true
	}
	nc.mountpoint = mountArr[0]

	ok := subscribe(nc.mountpoint, nc.id, nc.isMain, nc)
	if !ok {
		log.Warnf("%s subscribe %s failed", nc.id, nc.mountpoint)
		nc.conn.Close()
		return
	}
	log.Infof("%s isMain %t subscribe %s successfully", nc.id, nc.isMain, nc.mountpoint)

	nc.isReady = true
	nc.writeBuffer <- IcyOk
}

func (nc *ntripClient) parseBody() {
	buf := nc.buf[:]
	if nc.isMain {
		sendNmea(nc.mountpoint, buf)
	}
	nc.buf = nc.buf[0:0]
}

func (nc *ntripClient) callClose() {
	nc.once.Do(func() {
		nc.closeChan <- struct{}{}
	})
}

func (nc *ntripClient) sendData(buf []byte) {
	_, err := nc.conn.Write(buf)
	if err != nil {
		log.Errorf("%s write data error: %s", nc.id, err.Error())
		nc.callClose()
		return
	}
}

func (nc *ntripClient) IsMain() bool {
	return nc.isMain
}

func (nc *ntripClient) OnData(msg SubMsg) {
	if nc.isClosed {
		return
	}

	if msg.action == SendAction {
		nc.writeBuffer <- msg.data
	}

	if msg.action == TickAction {
		log.Infof("%s do tick action", nc.id)
		nc.conn.Close()
	}
}
