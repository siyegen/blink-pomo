package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func jsonEndpoint(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log.Print(LogPrefix, req.URL)
		res.Header().Set("Content-Type", "application/json")
		handler(res, req)
	})
}

func (b *BlinkApp) CreateChrono(res http.ResponseWriter, req *http.Request) {
	pom := NewPom()
	b.StorePom(pom)

	logLine(fmt.Sprintf("Creating Pom %s", pom.id))
	// go pom.StartTimer()
	res.Write(pom.ToJSON())
}

func (b *BlinkApp) StopPom(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	logLine(fmt.Sprintf("Stoping pom for %s", vars["id"]))

	pom, ok := b.currentPoms[vars["id"]]
	if !ok {
		logLine("No pom :(")
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(`{"error": "No Pom Found"}`))
		return
	}
	pom.StopTimer()
}

func (b *BlinkApp) GetPom(res http.ResponseWriter, req *http.Request) {
	logLine("Loading exsisting Pom")
	vars := mux.Vars(req)
	pom, ok := b.currentPoms[vars["id"]]
	if !ok {
		logLine("No pom :(")
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(`{"error": "No Pom Found"}`))
		return
	}
	logLine(fmt.Sprintf("Pom-%s => seconds %d", vars["id"], pom.seconds))
	res.Write(pom.ToJSON())
}
