package main

import (
	"net/http"
	"os"
	"strings"

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

	var (
		mongoDBURI = os.Getenv("MONGO_URL")
		httpPort   = os.Getenv("HTTP_PORT")
		APIBase    = "/" + arena.APIBase
	)

	zap.L().Debug("Starting backend")
	a := arena.NewArenaServer(mongoDBURI)
	zap.L().Debug("Connected to MongoDB")

	router := mux.NewRouter()
	router.HandleFunc(strings.Join([]string{APIBase, "managers", "{id}"}, "/"), a.GetManager).Methods("GET")
	router.HandleFunc(strings.Join([]string{APIBase, "managers"}, "/"), a.NewManager).Methods("POST")
	router.HandleFunc(strings.Join([]string{APIBase, "managers", "{id}"}, "/"), a.DeleteManager).Methods("DELETE")
	router.HandleFunc(strings.Join([]string{APIBase, "managers", "{id}", "gladiators"}, "/"), a.GetGladiators).Methods("GET")
	router.HandleFunc(strings.Join([]string{APIBase, "managers", "{id}", "gladiators", "{id}"}, "/"), a.GetGladiator).Methods("GET")
	router.HandleFunc(strings.Join([]string{APIBase, "managers", "{id}", "gladiators"}, "/"), a.NewGladiator).Methods("POST")
	router.HandleFunc(strings.Join([]string{APIBase, "managers", "{id}", "gladiators", "{id}", "fight"}, "/"), a.FightGladiator).Methods("POST")
	router.HandleFunc(strings.Join([]string{APIBase, "managers", "{id}", "gladiators", "{id}"}, "/"), a.DeleteGladiator).Methods("DELETE")

	zap.L().Fatal("Error when serving",
		zap.Error(http.ListenAndServe(":"+httpPort, router)),
	)
}
