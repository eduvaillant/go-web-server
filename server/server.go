package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetAllPlayersScores() map[string]int
}

type PlayerServer struct {
	store PlayerStore	
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.store = store
	
	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playerHandler))

	p.Handler = router

	return p
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func(p *PlayerServer) playerHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	switch r.Method {
	case http.MethodPost:
			p.processWin(w, player)
	case http.MethodGet:
		if player != "" {
			p.showScore(w, player)
			break
		} 
		p.listAllPlayers(w)
	}
}

func(p *PlayerServer) processWin(w http.ResponseWriter, player string) {
  p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)
	w.Header().Set("Content-Type", "application/json")	

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}	
	
	json.NewEncoder(w).Encode(map[string]int{"score": score})
}

func(p *PlayerServer) listAllPlayers(w http.ResponseWriter) {
	scores := p.store.GetAllPlayersScores()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scores)
}