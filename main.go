package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/iamcathal/lofilibrarian/dtos"
	"github.com/iamcathal/lofilibrarian/endpoints"
	"github.com/iamcathal/lofilibrarian/goodreads"
	"github.com/iamcathal/lofilibrarian/openlibrary"
	"github.com/iamcathal/lofilibrarian/rabbitmq"
	"github.com/iamcathal/lofilibrarian/util"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	ApplicationStartUpTime time.Time
	InfluxDBClient         influxdb2.Client
)

func initConfig() dtos.AppConfig {
	return dtos.AppConfig{
		ApplicationStartUpTime: time.Now(),
	}
}

func initInfluxClient() {
	client := influxdb2.NewClientWithOptions(
		os.Getenv("INFLUXDB_URL"),
		os.Getenv("LOFI_BUCKET_TOKEN"),
		influxdb2.DefaultOptions().SetBatchSize(10))
	InfluxDBClient = client
}

func main() {
	godotenv.Load()

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

	rabbitmq.InitConnection()

	if os.Getenv("INFLUX_ENABLE") != "" {
		initInfluxClient()
		endpoints.InitInfluxClient(InfluxDBClient)
	}

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
