package main

import (
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
	return &BlinkApp{make(map[string]*Pom)}
}

func (b *BlinkApp) StorePom(pom *Pom) {
	b.currentPoms[pom.id] = pom
}

func main() {
	fmt.Println("blink-pomo: Pretty lights and productivity")

	app := NewBlinkApp()
	r := mux.NewRouter()

	r.HandleFunc("/status/{id}", func(res http.ResponseWriter, req *http.Request) {
		id := mux.Vars(req)["id"]
		pom, ok := app.currentPoms[id]
		if !ok {
			logLine(fmt.Sprintf("No pom: %s", id))
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte(`{"error": "No Pom Found"}`))
			res.Write([]byte("#000000\n"))
			return
		}
		logLine("checking status")
		if pom.seconds > 45 {
			res.Write([]byte("#00FF00\n"))
			return
		} else if pom.seconds > 10 {
			res.Write([]byte("#FF0000\n"))
			return
		}
		res.Write([]byte("#FFFFFF\n"))
	})

	// API Endpoints
	r.HandleFunc("/pom", jsonEndpoint(app.CreateChrono)).Methods("POST")
	// r.HandleFunc("/pom/{id}", jsonEndpoint(app.CreatePom)).Methods("PUT")
	r.HandleFunc("/pom/stop/{id}", jsonEndpoint(app.StopPom)).Methods("POST")

	r.HandleFunc("/pom/start/{id}", func(res http.ResponseWriter, req *http.Request) {
		logLine("Starting pom for exsisting timer")
		vars := mux.Vars(req)
		res.Write([]byte(fmt.Sprintf("endpoint: /pom%s", vars["id"])))
	}).Methods("POST")

	r.HandleFunc("/pom/{id}", jsonEndpoint(app.GetPom)).Methods("GET")

	// Static calls
	r.HandleFunc("/chrono/{id}", func(res http.ResponseWriter, req *http.Request) {
		logLine(fmt.Sprintf("loading app for %s", mux.Vars(req)["id"]))
		http.ServeFile(res, req, "assets/chrono")
	})
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("assets")))
	log.Fatal(http.ListenAndServe(":9913", r))
}
