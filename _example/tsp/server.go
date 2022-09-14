package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/arl/evolve"
	"github.com/arl/evolve/engine"
)

// a server shows the current state of the TSP via an HTTP server.
type server struct {
	solutions chan []point // channel type is temporary (should probably be a channel of Individual)
}

func newServer() *server {
	return &server{
		solutions: make(chan []point),
	}
}

func (s *server) serve(host string) {
	log.Printf("Server started, point your browser to http://%s\n", host)

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", s.ws)

	go http.ListenAndServe(host, nil)
}

func genotypeToPhenotype(genotype []int, cities []point) []point {
	pts := make([]point, len(genotype))
	for i := 0; i < len(pts); i++ {
		pts[i] = cities[genotype[i]]
	}

	return pts
}

func (s *server) start() error {
	// cities := randomPath(numCities)
	cities := berlin52

	plotSolution := func(sol []int) {
		s.solutions <- genotypeToPhenotype(sol, cities)
	}

	obs := engine.ObserverFunc(func(stats *evolve.PopulationStats[[]int]) {
		if stats.Generation%plotEach == 0 {
			plotSolution(stats.Best)
			fmt.Printf("[%d]: distance: %v\n", stats.Generation, stats.BestFitness)
		}
	})

	pop, err := runTSP(cities, obs)
	if err != nil {
		return fmt.Errorf("runTSP failed: %v", err)
	}

	log.Printf("quitting")
	plotSolution(pop.Candidates[0])

	return nil
}

func randomPath(length int) []point {
	path := make([]point, length)
	for i := range path {
		path[i].X = rand.Intn(xmax)
		path[i].Y = rand.Intn(ymax)
	}

	return path
}

func (s *server) ws(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	fmt.Println("ws start")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	if err = s.runTSP(ws); err != nil {
		log.Println(err)
	}

	fmt.Println("ws end")
}

type message struct {
	Type    string
	Payload interface{}
}

func sendInitMessage(conn *websocket.Conn, w, h int) error {
	type initMessage struct {
		Width, Height int
	}

	m := message{
		Type:    "init",
		Payload: initMessage{Width: w, Height: h},
	}

	return conn.WriteJSON(m)
}

func sendSolutionMessage(conn *websocket.Conn, solution []point) error {
	type solutionMessage struct {
		Solution []point
	}

	m := message{
		Type:    "solution",
		Payload: solutionMessage{Solution: solution},
	}

	return conn.WriteJSON(m)
}

func (s *server) runTSP(conn *websocket.Conn) error {
	if err := sendInitMessage(conn, xmax, ymax); err != nil {
		return fmt.Errorf("sending initTSP message: %v", err)
	}

	var m message
	if err := conn.ReadJSON(&m); err != nil {
		return err
	}

	if m.Type != "start" {
		return fmt.Errorf("expected type=start message, got type=%v", m.Type)
	}

	for sol := range s.solutions {
		if err := sendSolutionMessage(conn, sol); err != nil {
			return fmt.Errorf("sending initTSP message: %v", err)
		}
	}

	return nil
}
