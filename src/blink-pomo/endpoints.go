package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type pomAction struct {
	State string `json:"state"`
}

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

func (b *BlinkApp) UpdatePom(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	logLine(fmt.Sprintf("Updating pom for %s", vars["id"]))

	pom, ok := b.currentPoms[vars["id"]]
	if !ok {
		logLine("No pom :(")
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(`{"error": "No Pom Found"}`))
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logLine("No valid body")
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(`{"error": "Couldn't ready request"}`))
		return
	}

	var action pomAction
	err = json.Unmarshal(body, &action)
	if err != nil {
		logLine("No action")
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(`{"error": "Invalid Request"}`))
		return
	}
	if action.State == "starting" {
		pom.StartTimer()
		return
	} else if action.State == "stopping" {
		pom.StopTimer()
		return
	}

	logLine(fmt.Sprintf("Invalid action %s", action.State))
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(`{"error": "Invalid Request"}`))
	return
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
