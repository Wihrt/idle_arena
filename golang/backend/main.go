package main

import (
	"log"
	"net/http"

	"github.com/wihrt/idle_arena/arena/arena"
	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

func main() {

	var mongoDBURI = "mongodb://localhost:27017"

	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)

	a := arena.NewArena(mongoDBURI)

	router := mux.NewRouter()
	router.HandleFunc("/gladiator/", a.GetGladiator).Methods("GET")
	router.HandleFunc("/gladiator/new", a.CreateGladiator).Methods("POST")
	router.HandleFunc("/gladiator/fight", a.FightGladiator).Methods("POST")
	router.HandleFunc("/gladiator/kill", a.KillGladiator).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":5000", router))
}
