package main

import (
	"log"
	"net/http"
)

const healthAPIPort = ":5050"

func startHealthCheckAPI() {

	privateMux := http.NewServeMux()
	privateMux.HandleFunc("/healthz", readinessProbe)

	log.Fatal(http.ListenAndServe(healthAPIPort, privateMux))
}

func readinessProbe(w http.ResponseWriter, r *http.Request) {
	if instanceFactory.Initialized() {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	} else {
		w.WriteHeader(500)
		w.Write([]byte("not ready"))
	}
}
