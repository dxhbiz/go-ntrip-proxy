package proxy

import (
	"fmt"
	"net"

	"github.com/dxhbiz/go-ntrip-proxy/pkg/config"
	"github.com/dxhbiz/go-ntrip-proxy/pkg/kit/log"
)

const (
	ReadChanLength    = 64
	WirteChanLength   = 64
	ReadDefaultLength = 1024
)

// init ntrip server
func initNtripServer() {
	cfg := config.GetConfig()

	svrCfg := cfg.Server

	addr := fmt.Sprintf("%s:%d", svrCfg.Host, svrCfg.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("ntrip server start error %s", err.Error())
	}
	log.Infof("ntrip server listening on %s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("ntrip accept error %s", err.Error())
			return
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	client := newNtripClient(conn)
	go client.run()
}
