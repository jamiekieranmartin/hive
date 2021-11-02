package network

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shirou/gopsutil/net"
)

type Network struct {
	Netstat    []net.IOCountersStat `json:"netstat"`
	Interfaces []net.InterfaceStat  `json:"interfaces"`
}

func get() *Network {
	u, err := net.IOCounters(true)

	if err != nil {
		fmt.Printf("Error getting netstat info: %v\n", err)
	}

	i, err := net.Interfaces()

	if err != nil {
		fmt.Printf("Error getting interfaces info: %v\n", err)
	}

	return &Network{
		u,
		i,
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	d := get()

	json.NewEncoder(w).Encode(d)
}
