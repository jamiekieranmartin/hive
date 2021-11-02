package disk

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/disk"
)

type Disk struct {
	Block      *ghw.BlockInfo  `json:"block"`
	Statistics *disk.UsageStat `json:"statistics"`
}

func getDisk() *Disk {
	b, err := ghw.Block()

	if err != nil {
		fmt.Printf("Error getting block info: %v\n", err)
	}

	u, err := disk.Usage("/")

	if err != nil {
		fmt.Printf("Error getting disk info: %v\n", err)
	}

	return &Disk{
		b,
		u,
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	d := getDisk()

	json.NewEncoder(w).Encode(d)
}
