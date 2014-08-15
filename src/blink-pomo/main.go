package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const LogPrefix = "[BlinkApp]"

type ActivePoms map[string]*Pom

// Represents our server
type BlinkApp struct {
	currentPoms ActivePoms
}

func NewBlinkApp() *BlinkApp {
	return &BlinkApp{make(ActivePoms)}
}

func (a ActivePoms) StorePom(pom *Pom) {
	a[pom.id] = pom
}
func (a ActivePoms) RetrievePom(id string) (*Pom, error) {
	pom, ok := a[id]
	if !ok {
		return nil, errors.New("No Pom by that id")
	}
	return pom, nil
}

func main() {
	fmt.Println("blink-pomo: Pretty lights and productivity")

	app := NewBlinkApp()
	r := mux.NewRouter()

	r.HandleFunc("/status/{id}", app.PomStatus).Methods("GET")

	// API Endpoints
	r.HandleFunc("/pom", jsonEndpoint(app.CreateChrono)).Methods("POST")
	r.HandleFunc("/pom/{id}", jsonEndpoint(app.UpdatePom)).Methods("PUT")
	r.HandleFunc("/pom/{id}", jsonEndpoint(app.GetPom)).Methods("GET")

	// Static calls
	r.HandleFunc("/chrono/{id}", func(res http.ResponseWriter, req *http.Request) {
		logLine(fmt.Sprintf("loading app for %s", mux.Vars(req)["id"]))
		http.ServeFile(res, req, "assets/chrono")
	})
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("assets")))
	log.Fatal(http.ListenAndServe(":9913", r))
}
