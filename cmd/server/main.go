package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/gauge", GaugeHandler)
	mux.HandleFunc("/update/counter", CounterHandler)

	http.ListenAndServe("localhost:8080", mux)
}

func GaugeHandler(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodPost
	r.Header.Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	body := "Query parameters ===============\r\n"
	for k, v := range r.URL.Query() {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	w.Write([]byte(body))
}

func CounterHandler(w http.ResponseWriter, r *http.Request) {
	r.Method = http.MethodPost
	w.Write([]byte("test\r\n"))
}
