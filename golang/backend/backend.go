package main

import (
	"net/http"
	"os"

	"github.com/wihrt/idle_arena/arena"
	"github.com/wihrt/idle_arena/logging"
	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

func init() {

	cfg := logging.GetConfig()
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}

func main() {

	var mongoDBURI = os.Getenv("MONGO_URL")

	zap.L().Debug("Starting backend")
	a := arena.NewArenaServer(mongoDBURI)
	zap.L().Debug("Connected to MongoDB")

	router := mux.NewRouter()
	router.HandleFunc("/"+arena.APIBase+"/managers/{id}", a.GetManager).Methods("GET")
	router.HandleFunc("/"+arena.APIBase+"/managers", a.NewManager).Methods("POST")
	router.HandleFunc("/"+arena.APIBase+"/managers/{id}", a.DeleteManager).Methods("DELETE")
	router.HandleFunc("/"+arena.APIBase+"/managers/{id}/gladiators", a.GetGladiators).Methods("GET")
	router.HandleFunc("/"+arena.APIBase+"/managers/{id}/gladiators/{id}", a.GetGladiator).Methods("GET")
	router.HandleFunc("/"+arena.APIBase+"/managers/{id}/gladiators", a.NewGladiator).Methods("POST")
	router.HandleFunc("/"+arena.APIBase+"/managers/{id}/gladiators/{id}/fight", a.FightGladiator).Methods("POST")
	router.HandleFunc("/"+arena.APIBase+"/managers/{id}/gladiators/{id}", a.DeleteGladiator).Methods("DELETE")

	zap.L().Fatal("Error when serving",
		zap.Error(http.ListenAndServe(":5000", router)),
	)
}
