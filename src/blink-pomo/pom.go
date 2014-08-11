package main

import (
	"encoding/json"
	"log"
	"time"
)

type pomResponse struct {
	UUID      string `json:"uuid"`
	StartTime int64  `json:"start_time"`
}

type Pom struct {
	timer     *time.Timer
	ticker    *time.Ticker
	flag      string
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
	p.flag = "started"
	for t := range p.ticker.C {
		log.Print(t)
		p.seconds += 5
	}
}

func (p *Pom) StopTimer() {
	p.ticker.Stop()
	p.flag = "stopped"
}
