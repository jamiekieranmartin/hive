package service

import (
	"log"
	"net"
	"net/http"

	"github.com/gobwas/ws"
)

// Node object
type Node struct {
	Header  ws.HandshakeHeaderHTTP
	Upgrade *ws.Upgrader
	Handler func(ctx Context) error
}

// Start service
func (s *Node) Start(ip string) error {
	ln, err := net.Listen("tcp", ip)
	if err != nil {
		return err
	}

	if s.Header == nil {
		s.Header = ws.HandshakeHeaderHTTP(http.Header{})
	}

	if s.Upgrade == nil {
		s.Upgrade = &ws.Upgrader{
			OnHeader: func(key, value []byte) error {
				return nil
			},
		}
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		log.Println("new connection")

		_, err = s.Upgrade.Upgrade(conn)
		if err != nil {
			return err
		}

		defer conn.Close()

		go func() {
			err := s.Handler(Context{conn})
			if err != nil {
				log.Println(err)
			}
		}()
	}
}
