package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHandler)

	http.ListenAndServe("localhost:8080", mux)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}
