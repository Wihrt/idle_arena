package main

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/wihrt/idle_arena/dnd"
	"github.com/wihrt/idle_arena/logging"
	"go.uber.org/zap"
)

func init() {
	cfg := logging.GetConfig()
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}

func main() {

	var (
		mongoURI   = os.Getenv("MONGO_URL")
		dndAPI     = os.Getenv("DND_API")
		timeoutStr = os.Getenv("TIMEOUT")
		indexes    = os.Getenv("INDEXES")
		wg         sync.WaitGroup
	)

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		zap.L().Fatal("Cannot parse duration",
			zap.String("duration", timeoutStr),
			zap.Error(err),
		)
	}

	loader := dnd.NewLoader(dndAPI, mongoURI, timeout)

	for _, indexName := range strings.Split(indexes, ",") {
		wg.Add(1)
		go loader.LoadIndex(indexName, &wg)
	}
	wg.Wait()
}
