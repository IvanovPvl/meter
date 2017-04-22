// Package server
package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ivanovpvl/meter/agent/df"
)

const defaultPort int = 8888

func init() {
	log.SetOutput(os.Stdout)
}

func handler(w http.ResponseWriter, r *http.Request) {
	res, err := df.Exec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Run server
func Run() {
	port := flag.Int("port", defaultPort, "Port for listening")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	log.Printf("Agent running on %d port.\n", *port)
	http.ListenAndServe(addr, mux)
}
