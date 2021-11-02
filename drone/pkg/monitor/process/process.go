package process

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shirou/gopsutil/process"
)

type Process struct {
	Pid    int32   `json:"pid"`
	Name   string  `json:"name"`
	Exe    string  `json:"exe"`
	CPU    float64 `json:"cpu"`
	Memory float32 `json:"memory"`
}

func getProcess(p *process.Process) *Process {
	n, _ := p.Name()
	e, _ := p.Exe()
	c, _ := p.CPUPercent()
	m, _ := p.MemoryPercent()

	return &Process{
		p.Pid,
		n,
		e,
		c,
		m,
	}
}

func getProcesses() []*Process {
	p, err := process.Processes()

	if err != nil {
		fmt.Printf("Error getting process info: %v\n", err)
	}

	var list []*Process

	for _, pr := range p {
		process := getProcess(pr)
		list = append(list, process)
	}

	return list
}

func Handler(w http.ResponseWriter, r *http.Request) {
	list := getProcesses()

	json.NewEncoder(w).Encode(list)
}
