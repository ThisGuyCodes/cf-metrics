package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

var metrics = make(map[string]func() float64)
var metricMutex = sync.Mutex{}

type metric struct {
	name string
	tags map[string]string
}

func (m metric) String() string {
	return fmt.Sprintf("%s{%v}", m.name, m.tags)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	metricMutex.Lock()
	defer metricMutex.Unlock()
	for metric, value := range metrics {
		w.Write([]byte(fmt.Sprintf("%s %v\n", metric, value())))
	}
}

var machine_id = "5"

func main() {
	metrics["my_app_memory{pid=\"12345\"}"] = func() float64 { return 1234 }
	metrics[fmt.Sprintf("process_count{machine_id=%q}", machine_id)] = get_process_count
	metrics[fmt.Sprintf("log_file_size{machine_id=%q}", machine_id)] = log_file_size

	http.HandleFunc("/metrics", metricsHandler)
	http.ListenAndServe(":10000", nil)
}

func get_process_count() float64 {
	return 4
}

func log_file_size() float64 {
	f, err := os.Open("log.txt")
	defer f.Close()
	if err != nil {
		// log something that sets of alarms
		return 0
	}

	stat, err := f.Stat()
	if err != nil {
		// log something that sets of alarms
		return 0
	}

	return float64(stat.Size())
}
