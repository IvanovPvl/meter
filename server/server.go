package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"gopkg.in/yaml.v2"
	"github.com/ivanovpvl/meter/agent/df"
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
			resp, _ := http.Get(host)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				// TODO: handle error
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
			fmt.Println(err)
		}

		hostResp := hostResponse{result.Host, dfRes}
		data = append(data, hostResp)
	}

	resp, err := json.Marshal(data)
	if err != nil {
		// TODO
	}

	w.Write(resp)
}

// Run monitoring server
func Run() {
	cnf, err := parseConfig()
	if err != nil {
		fmt.Println(err)
	}
	addr := fmt.Sprintf(":%d", cnf.Port)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/df", dfHandler)
	http.ListenAndServe(addr, mux)
}
