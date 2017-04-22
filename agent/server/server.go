// Package server
package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ivanovpvl/meter/agent/df"
	"net/http"
)

const defaultPort int = 8888

func handler(w http.ResponseWriter, r *http.Request) {
	res, err := df.Exec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Start server
func Start() {
	port := flag.Int("port", defaultPort, "Port for listening")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	http.ListenAndServe(addr, mux)
}
