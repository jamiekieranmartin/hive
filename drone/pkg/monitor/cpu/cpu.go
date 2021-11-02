package cpu

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jaypipes/ghw"
)

func get() *ghw.CPUInfo {
	c, err := ghw.CPU()

	if err != nil {
		fmt.Printf("Error getting cpu info: %v\n", err)
	}

	return c
}

func Handler(w http.ResponseWriter, r *http.Request) {
	c := get()

	json.NewEncoder(w).Encode(c)
}
