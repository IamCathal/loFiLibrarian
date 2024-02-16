package influxdb

import (
	"context"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

var (
	influxClient influxdb2.Client
)

func IsInfluxDBEnabled() bool {
	return os.Getenv("OPT_INFLUX_ENABLE") == "true"
}

func InitInfluxClient() {
	newClient := influxdb2.NewClientWithOptions(
		os.Getenv("INFLUXDB_URL"),
		os.Getenv("LOFI_BUCKET_TOKEN"),
		influxdb2.DefaultOptions().SetBatchSize(10))
	influxClient = newClient
}

func WriteMetricPoint(ip string) {
	writeAPI := influxClient.WriteAPIBlocking(os.Getenv("ORG"), os.Getenv("LOFI_BUCKET"))
	point := influxdb2.NewPointWithMeasurement("clientIPLog").
		AddTag("clientIP", ip).
		AddField("service", "lofilibrarian").
		SetTime(time.Now())
	writeAPI.WritePoint(context.Background(), point)
}
