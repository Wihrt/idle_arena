package main

import (
	"net/http"
	"os"

	"github.com/wihrt/idle_arena/arena/arena"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gorilla/mux"
)

const APIBase = "/api/v1"

func init() {

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		DisableStacktrace: false,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}

func main() {

	var mongoDBURI = os.Getenv("MONGO_URL")

	zap.L().Debug("Starting backend")
	a := arena.NewArena(mongoDBURI)
	zap.L().Debug("Connected to MongoDB")

	router := mux.NewRouter()
	router.HandleFunc(APIBase+"/managers/{id}", a.GetManager).Methods("GET")
	router.HandleFunc(APIBase+"/managers/{id}", a.DeleteManager).Methods("DELETE")
	router.HandleFunc(APIBase+"/managers/{id}/gladiators", a.GetGladiators).Methods("GET")
	router.HandleFunc(APIBase+"/managers/{id}/gladiators/{id}", a.GetGladiator).Methods("GET")
	router.HandleFunc(APIBase+"/managers/{id}/gladiators", a.NewGladiator).Methods("POST")
	router.HandleFunc(APIBase+"/managers/{id}/gladiators/{id}/fight", a.FightGladiator).Methods("POST")
	router.HandleFunc(APIBase+"/managers/{id}/gladiators/{id}", a.DeleteGladiator).Methods("DELETE")

	zap.L().Fatal("Error when serving",
		zap.Error(http.ListenAndServe(":5000", router)),
	)
}
