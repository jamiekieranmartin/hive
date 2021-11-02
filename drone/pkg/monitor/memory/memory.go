package memory

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/mem"
)

type Memory struct {
	Hardware   *ghw.MemoryInfo        `json:"hardware"`
	Statistics *mem.VirtualMemoryStat `json:"statistics"`
}

func getMemory() *Memory {
	m, err := ghw.Memory()

	if err != nil {
		fmt.Printf("Error getting memory info: %v\n", err)
	}

	v, err := mem.VirtualMemory()

	if err != nil {
		fmt.Printf("Error getting memory statistics: %v\n", err)
	}

	return &Memory{
		m,
		v,
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	m := getMemory()

	json.NewEncoder(w).Encode(m)
}
