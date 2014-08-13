package main

import (
	"encoding/json"
	"log"
	"time"
)

type pomResponse struct {
	UUID      string   `json:"uuid"`
	StartTime int64    `json:"start_time"`
	State     PomState `json:"state"`
	Seconds   int      `json:"seconds"`
}

type Pom struct {
	timer     *time.Timer
	ticker    *time.Ticker
	state     PomState
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
	pomRes := pomResponse{p.id, p.startTime, p.state, p.seconds}
	jsonRes, _ := json.Marshal(pomRes)
	return jsonRes
}

func (p *Pom) StartTimer() {
	p.state = pomStart
	for t := range p.ticker.C {
		log.Print(t)
		p.seconds += 5
	}
}

func (p *Pom) StopTimer() {
	p.ticker.Stop()
	p.state = pomStop
}
