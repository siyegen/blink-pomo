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

const LogPrefix = "[BlinkApp]"

func jsonEndpoint(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log.Print(LogPrefix, req.URL)
		res.Header().Set("Content-Type", "application/json")
		handler(res, req)
	})
}

type BlinkApp struct {
	currentPoms map[string]*Pom
}

func NewBlinkApp() *BlinkApp {
	return &BlinkApp{make(map[string]*Pom)}
}

func logLine(msg string) {
	log.Printf("%s %s\n", LogPrefix, msg)
}

func (b *BlinkApp) StartPom(res http.ResponseWriter, req *http.Request) {
	pom := NewPom()
	b.StorePom(pom)

	logLine(fmt.Sprintf("Creating Pom %s", pom.id))
	go pom.StartTimer()
	res.Write(pom.ToJSON())
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

func (b *BlinkApp) StorePom(pom *Pom) {
	b.currentPoms[pom.id] = pom
}

func newUUID() string {
	h := md5.New()
	b := make([]byte, 16)
	rand.Read(b)
	io.WriteString(h, string(b))
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

type pomResponse struct {
	UUID      string `json:"uuid"`
	StartTime int64  `json:"start_time"`
}

type Pom struct {
	timer     *time.Timer
	ticker    *time.Ticker
	startTime int64
	seconds   int
	id        string
}

func NewPom() *Pom {
	return &Pom{
		timer:     time.NewTimer(25 * time.Minute),
		ticker:    time.NewTicker(5 * time.Second),
		startTime: time.Now().Unix(),
		seconds:   0,
		id:        newUUID(),
	}
}

func (p *Pom) ToJSON() []byte {
	pomRes := pomResponse{p.id, p.startTime}
	jsonRes, _ := json.Marshal(pomRes)
	return jsonRes
}

func (p *Pom) StartTimer() {
	for t := range p.ticker.C {
		log.Print(t)
		p.seconds += 5
	}
}

func (p *Pom) StopTimer() {
	p.ticker.Stop()
}

func main() {
	fmt.Println("blink-pomo: Pretty lights and productivity")

	app := NewBlinkApp()
	r := mux.NewRouter()

	r.HandleFunc("/status/{id}", func(res http.ResponseWriter, req *http.Request) {
		id := mux.Vars(req)["id"]
		pom := app.currentPoms[id]
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

	//consider grouping under api
	r.HandleFunc("/pom", jsonEndpoint(app.StartPom)).Methods("POST")

	r.HandleFunc("/pom/stop/{id}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		logLine(fmt.Sprintf("Stoping pom for %s", vars["id"]))

		pom, ok := app.currentPoms[vars["id"]]
		if !ok {
			logLine("No pom :(")
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte(`{"error": "No Pom Found"}`))
			return
		}
		pom.StopTimer()
	}).Methods("POST")

	r.HandleFunc("/pom/start/{id}", func(res http.ResponseWriter, req *http.Request) {
		logLine("Starting pom for exsisting timer")
		vars := mux.Vars(req)
		res.Write([]byte(fmt.Sprintf("endpoint: /pom%s", vars["id"])))
	}).Methods("POST")

	r.HandleFunc("/pom/{id}", jsonEndpoint(app.GetPom)).Methods("GET")

	r.HandleFunc("/p/{id}", func(res http.ResponseWriter, req *http.Request) {
		logLine(fmt.Sprintf("loading app for %s", mux.Vars(req)["id"]))
		http.ServeFile(res, req, "assets")
	})
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("assets")))
	http.ListenAndServe(":9913", r)
}
