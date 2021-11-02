package host

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/host"
)

type Host struct {
	Host    *host.InfoStat   `json:"host"`
	Product *ghw.ProductInfo `json:"product"`
}

func getHost() *Host {
	h, err := host.Info()

	if err != nil {
		fmt.Printf("Error getting host info: %v\n", err)
	}

	p, err := ghw.Product()

	if err != nil {
		fmt.Printf("Error getting Product info: %v\n", err)
	}

	return &Host{
		h,
		p,
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	h := getHost()

	json.NewEncoder(w).Encode(h)
}
