package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"time"
)

type BlinkApp struct {
	currentPoms map[string]*Pom
}

func NewBlinkApp() *BlinkApp {
	return &BlinkApp{make(map[string]*Pom)}
}

func newUUID() string {
	h := md5.New()
	b := make([]byte, 16)
	rand.Read(b)
	io.WriteString(h, string(b))
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

type pomResponse struct {
	UUID string `json:"uuid"`
}

type Pom struct {
	timer     *time.Timer
	ticker    *time.Ticker
	startTime int64
	seconds   int
}

func NewPom() *Pom {
	return &Pom{
		time.NewTimer(25 * time.Minute),
		time.NewTicker(5 * time.Second),
		time.Now().Unix(),
		0,
	}
}

func main() {
	fmt.Println("blink-pomo: Pretty lights and productivity")

	app := NewBlinkApp()
	r := mux.NewRouter()

	r.HandleFunc("/pom", func(res http.ResponseWriter, req *http.Request) {
		log.Print("Creating new Pom")

		pom := NewPom()
		go func() {
			for t := range pom.ticker.C {
				log.Print(t)
				pom.seconds += 5
			}
		}()
		pomRes := pomResponse{newUUID()}
		app.currentPoms[pomRes.UUID] = pom
		jsonRes, _ := json.Marshal(pomRes)
		res.Write(jsonRes)
	}).Methods("POST")

	r.HandleFunc("/pom/{id}", func(res http.ResponseWriter, req *http.Request) {
		log.Print("Starting pom for exsisting timer")
		vars := mux.Vars(req)
		res.Write([]byte(fmt.Sprintf("endpoint: /pom%s", vars["id"])))
	}).Methods("POST")

	r.HandleFunc("/pom/{id}", func(res http.ResponseWriter, req *http.Request) {
		log.Print("Loading exsisting Pom")
		vars := mux.Vars(req)
		pom, ok := app.currentPoms[vars["id"]]
		if !ok {
			log.Print("No pom :(")
			res.Write([]byte("No pom"))
		}
		log.Printf("Pom-%s => seconds %d", vars["id"], pom.seconds)
	}).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("assets")))
	http.ListenAndServe(":9913", r)
}
