package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"runtime"

	"github.com/gobwas/httphead"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type user struct {
	ID int `json:"id"`
}

func main() {
	ln, err := net.Listen("tcp", "localhost:8080")
	check(err)

	// Prepare handshake header writer from http.Header mapping.
	header := ws.HandshakeHeaderHTTP(http.Header{
		"X-Go-Version": []string{runtime.Version()},
	})

	u := ws.Upgrader{
		OnHost: func(host []byte) error {
			if string(host) == "localhost:8080" {
				return nil
			}
			return ws.RejectConnectionError(
				ws.RejectionStatus(403),
				ws.RejectionHeader(ws.HandshakeHeaderString(
					"X-Want-Host: localhost:8080\r\n",
				)),
			)
		},
		OnHeader: func(key, value []byte) error {
			if string(key) != "Cookie" {
				return nil
			}
			ok := httphead.ScanCookie(value, func(key, value []byte) bool {
				// Check session here or do some other stuff with cookies.
				// Maybe copy some values for future use.
				return true
			})
			if ok {
				return nil
			}
			return ws.RejectConnectionError(
				ws.RejectionReason("bad cookie"),
				ws.RejectionStatus(400),
			)
		},
		OnBeforeUpgrade: func() (ws.HandshakeHeader, error) {
			return header, nil
		},
	}
	for {
		conn, err := ln.Accept()
		check(err)

		_, err = u.Upgrade(conn)
		if err != nil {
			log.Printf("upgrade error: %s", err)
		}

		go func() error {
			defer conn.Close()

			r := wsutil.NewReader(conn, ws.StateServerSide)
			decoder := json.NewDecoder(r)

			for {
				hdr, err := r.NextFrame()
				if err != nil {
					return err
				}
				if hdr.OpCode == ws.OpClose {
					return io.EOF
				}

				var u user
				if err := decoder.Decode(&u); err != nil {
					return err
				}

				fmt.Println(u)
			}
		}()
	}

}
