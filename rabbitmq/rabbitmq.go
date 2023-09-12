package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/iamcathal/lofilibrarian/dtos"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger

	rabbitMqChan *amqp.Channel
	queueName    string
)

func SetLogger(newLogger *zap.Logger) {
	logger = newLogger
}

func InitConnection() (amqp.Queue, amqp.Channel) {
	fmt.Println(os.Getenv("OPT_RABBITMQ_USER"), os.Getenv("OPT_RABBITMQ_PASSWORD"), os.Getenv("OPT_RABBITMQ_URL"))

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s", os.Getenv("OPT_RABBITMQ_USER"), os.Getenv("OPT_RABBITMQ_PASSWORD"), os.Getenv("OPT_RABBITMQ_URL")))
	if err != nil {
		panic(err)
	}
	// defer conn.Close()
	logger.Sugar().Infof("Connected to rabbitMQ instance %s", os.Getenv("OPT_RABBITMQ_URL"))

	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	rabbitMqChan = channel
	// defer channel.Close()

	queue, err := channel.QueueDeclare(
		os.Getenv("OPT_RABBITMQ_QUEUE_NAME"), // name
		false,                                // durable
		false,                                // delete when unused
		false,                                // exclusive
		false,                                // no-wait
		nil,                                  // arguments
	)
	if err != nil {
		panic(err)
	}
	queueName = queue.Name

	err = channel.Qos(
		2,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		panic(err)
	}

	return queue, *channel
}

func SyncWriteBookLookup(breadcrumb dtos.BookBreadcrumb) error {
	breadCrumbJson, err := json.Marshal(breadcrumb)
	if err != nil {
		return errWithTrace(err)
	}
	return publish(breadCrumbJson)
}

func publish(msgJson []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	fmt.Printf("Publishing %s\n", string(msgJson))

	return rabbitMqChan.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/json",
			Body:        msgJson,
		},
	)
}

func errWithTrace(err error) error {
	return fmt.Errorf(err.Error(), string(debug.Stack()))
}

func IsRabbitMQEnabled() bool {
	return os.Getenv("OPT_RABBITMQ_ENABLE") == "true"
}
