package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/iamcathal/lofilibrarian/dtos"
	"github.com/iamcathal/lofilibrarian/endpoints"
	"github.com/iamcathal/lofilibrarian/goodreads"
	"go.uber.org/zap"
)

var (
	ApplicationStartUpTime time.Time
)

func initConfig() dtos.AppConfig {
	return dtos.AppConfig{
		ApplicationStartUpTime: time.Now(),
	}
}

func main() {
	logConfig := zap.NewProductionConfig()
	logConfig.OutputPaths = []string{"stdout", "logs/appLog.log"}
	globalLogFields := make(map[string]interface{})
	globalLogFields["service"] = "lofilibrarian"
	logConfig.InitialFields = globalLogFields

	logger, err := logConfig.Build()
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	appConfig := initConfig()
	endpoints.InitConfig(appConfig, logger)
	goodreads.SetLogger(logger)

	port := 2946

	router := endpoints.SetupRouter()

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + fmt.Sprint(port),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	logger.Sugar().Infof("Service requests on :%d", port)
	log.Fatal(srv.ListenAndServe())
}
