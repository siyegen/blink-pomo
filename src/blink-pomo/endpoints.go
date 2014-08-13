package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type PomState string

const (
	pomStart PomState = "started"
	pomStop  PomState = "stopped"
)

type pomAction struct {
	State PomState `json:"state"`
}

func jsonEndpoint(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log.Print(LogPrefix, req.URL)
		res.Header().Set("Content-Type", "application/json")
		handler(res, req)
	})
}

func (b *BlinkApp) PomStatus(res http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	pom, ok := b.currentPoms[id]
	if !ok {
		logLine(fmt.Sprintf("No pom: %s", id))
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(`{"error": "No Pom Found"}`))
		res.Write([]byte("#000000\n"))
		return
	}
	logLine("checking status")
	if pom.state == "started" {
		res.Write([]byte("#00FF00\n"))
		return
	} else if pom.state == "stopped" {
		res.Write([]byte("#FF0000\n"))
		return
	}
	res.Write([]byte("#FFFFFF\n"))
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
	if action.State == pomStart {
		pom.StartTimer()
		return
	} else if action.State == pomStop {
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
