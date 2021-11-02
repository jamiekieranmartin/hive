package main

import (
	"flag"
	"fmt"
	"hive/service"
	"log"
)

var (
	port = flag.String("port", "80", "port to run the service on")
	host = flag.String("host", "127.0.0.1", "host to run the service on")
)

func init() {
	flag.Parse()
}

func main() {
	node := service.Node{
		Handler: func(ctx service.Context) error {
			for {
				msg, op, err := ctx.Read()
				if err != nil {
					return err
				}

				fmt.Println(msg, op)

				err = ctx.Write(op, msg)
				if err != nil {
					return err
				}
			}
		},
	}

	addr := fmt.Sprintf("%s:%s", *host, *port)

	log.Println("listening on", addr)

	node.Start(addr)
}
