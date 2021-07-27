package main

import (
	"log"
	"net/http"
	"os"

	"github.com/wihrt/idle_arena/arena/arena"
	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

func main() {

	var mongoDBURI = os.Getenv("MONGO_URL")

	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)

	a := arena.NewArena(mongoDBURI)

	router := mux.NewRouter()
	router.HandleFunc("/gladiator/", a.GetGladiator).Methods("GET")
	router.HandleFunc("/gladiator/hire", a.HireGladiator).Methods("POST")
	router.HandleFunc("/gladiator/fight", a.FightGladiator).Methods("POST")
	router.HandleFunc("/gladiator/fire", a.FireGladiator).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":5000", router))
}
