package main

//
// import (
// 	"drone/pkg/monitor/cpu"
// 	"drone/pkg/monitor/disk"
// 	"drone/pkg/monitor/host"
// 	"drone/pkg/monitor/memory"
// 	"drone/pkg/monitor/network"
// 	"drone/pkg/monitor/process"
// 	"fmt"
// 	"log"
// 	"net/http"
// )

// func main() {

// 	http.HandleFunc("/cpu", cpu.Handler)
// 	http.HandleFunc("/disk", disk.Handler)
// 	http.HandleFunc("/host", host.Handler)
// 	http.HandleFunc("/memory", memory.Handler)
// 	http.HandleFunc("/network", network.Handler)
// 	http.HandleFunc("/process", process.Handler)

// 	fmt.Println("Listening on 5000")
// 	log.Fatal(http.ListenAndServe(":5000", nil))
// }

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/namsral/flag"
)

var addr = flag.String("addr", "localhost:80", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:

			p := getProcesses()

			js, err := json.Marshal(p)
			check(err)

			err = c.WriteMessage(websocket.TextMessage, js)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
