// Package server
package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/ivanovpvl/meter/agent/df"
	"gopkg.in/yaml.v2"
)

type config struct {
	Port    int
	Servers []string `yaml:"servers"`
}

type transfer struct {
	Host string
	Data []byte
}

type hostResponse struct {
	Host     string      `json:"host"`
	DiscFree []df.Result `json:"disk_free"`
}

var hosts []string

func init() {
	log.SetOutput(os.Stdout)
}

func parseConfig() (*config, error) {
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return nil, err
	}

	cnf := config{}
	err = yaml.Unmarshal(data, &cnf)
	hosts = cnf.Servers
	return &cnf, err
}

func dfHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(hosts))
	results := make(chan transfer)

	for _, host := range hosts {
		go func(host string) {
			resp, err := http.Get(host)
			if err != nil || resp.StatusCode != http.StatusOK {
				log.Println(host, err)
				waitGroup.Done()
				return
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(host, err)
				waitGroup.Done()
				return
			}

			result := transfer{host, body}
			results <- result
			waitGroup.Done()
		}(host)
	}

	go func() {
		waitGroup.Wait()
		close(results)
	}()

	processResults(results, w)
}

func processResults(results chan transfer, w http.ResponseWriter) {
	data := make([]hostResponse, 0)
	for result := range results {
		dfRes := []df.Result{}
		err := json.Unmarshal(result.Data, &dfRes)
		if err != nil {
			log.Println(result.Host, err)
			continue
		}

		hostResp := hostResponse{result.Host, dfRes}
		data = append(data, hostResp)
	}

	resp, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		errMap := make(map[string]error)
		errMap["message"] = err
		errResp, e := json.Marshal(errMap)
		if e != nil {
			log.Println(e)
			return
		}

		w.Write(errResp)
	}

	w.Write(resp)
}

// Run monitoring server
func Run() {
	cnf, err := parseConfig()
	if err != nil {
		log.Fatalln(err)
	}

	addr := fmt.Sprintf(":%d", cnf.Port)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/df", dfHandler)

	log.Printf("Server running on %d port.\n", cnf.Port)
	http.ListenAndServe(addr, mux)
}
