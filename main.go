package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/iamcathal/lofilibrarian/dtos"
	"github.com/iamcathal/lofilibrarian/endpoints"
	"github.com/iamcathal/lofilibrarian/goodreads"
	"github.com/iamcathal/lofilibrarian/influxdb"
	"github.com/iamcathal/lofilibrarian/openlibrary"
	"github.com/iamcathal/lofilibrarian/rabbitmq"
	"github.com/iamcathal/lofilibrarian/util"
	"github.com/joho/godotenv"
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
	godotenv.Load(".env", ".localenv")

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
	openlibrary.SetLogger(logger)
	util.SetLogger(logger)
	rabbitmq.SetLogger(logger)

	logger.Sugar().Infof("RabbitMQ enabled: %v", rabbitmq.IsRabbitMQEnabled())
	if rabbitmq.IsRabbitMQEnabled() {
		rabbitmq.InitConnection()
	}
	logger.Sugar().Infof("InfluxDB enabled: %v", influxdb.IsInfluxDBEnabled())
	if influxdb.IsInfluxDBEnabled() {
		influxdb.InitInfluxClient()
	}

	port := 2946

	router := endpoints.SetupRouter()

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + fmt.Sprint(port),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	logger.Sugar().Infof("Servicing requests on :%d", port)
	log.Fatal(srv.ListenAndServe())
}
