package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	// file web server
	mux.Handle("/", http.FileServer(http.Dir("web")))

	err := http.ListenAndServe(":7540", mux) // TODO: get port from env variables
	if err != nil {
		panic(err)
	}
}
