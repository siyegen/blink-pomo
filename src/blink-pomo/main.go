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

	// app := NewBlinkApp()
	// r := mux.NewRouter()

	// r.HandleFunc("/status/{id}", app.PomStatus).Methods("GET")

	// // API Endpoints
	// r.HandleFunc("/pom", jsonEndpoint(app.CreateChrono)).Methods("POST")
	// r.HandleFunc("/pom/{id}", jsonEndpoint(app.UpdatePom)).Methods("PUT")
	// r.HandleFunc("/pom/{id}", jsonEndpoint(app.GetPom)).Methods("GET")

	// // Static calls
	// r.HandleFunc("/chrono/{id}", func(res http.ResponseWriter, req *http.Request) {
	// 	logLine(fmt.Sprintf("loading app for %s", mux.Vars(req)["id"]))
	// 	http.ServeFile(res, req, "assets/chrono")
	// })
	// r.PathPrefix("/").Handler(http.FileServer(http.Dir("assets")))
	// log.Fatal(http.ListenAndServe(":9913", r))
	main2()
}

type Endpoint func(http.ResponseWriter, *http.Request)
type fancyServer struct {
	endpoints []Endpoint
	app       *BlinkApp
	router    *mux.Router
}

func main2() {

	app := NewBlinkApp()
	server := NewfancyServer(app)
	log.Fatal(server.Run())
}

func NewfancyServer(app *BlinkApp) *fancyServer {
	router := mux.NewRouter()
	pomStatusEndpoint := func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		res, err := app.CheckPomStatus(id)
		writeRes(res, err, w)
	}
	router.HandleFunc("/status/{id}", pomStatusEndpoint).Methods("GET")
	return &fancyServer{app: app, router: router}
}

func writeRes(res []byte, err error, w http.ResponseWriter) {
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (f *fancyServer) Run() error {
	return http.ListenAndServe(":9913", f.router)
}
func (b *BlinkApp) CheckPomStatus(id string) ([]byte, error) {
	pom, ok := b.currentPoms[id]
	if !ok {
		logLine(fmt.Sprintf("No pom: %s", id))
		return []byte("#000000\n"), errors.New("Not Found")
	}

	logLine("checking status")
	if pom.state == pomStart {
		return []byte("#00FF00\n"), nil
	} else if pom.state == pomStop {
		return []byte("#FF0000\n"), nil
	}
	return []byte("#FFFFFF\n"), nil
}
