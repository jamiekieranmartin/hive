package service

import (
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// Context object
type Context struct {
	conn net.Conn
}

func (ctx *Context) Read() ([]byte, ws.OpCode, error) {
	return wsutil.ReadClientData(ctx.conn)
}

func (ctx *Context) Write(op ws.OpCode, msg []byte) error {
	return wsutil.WriteServerMessage(ctx.conn, op, msg)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
