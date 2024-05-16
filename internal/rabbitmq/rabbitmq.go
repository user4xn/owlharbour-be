package rabbitmq

import (
	"fmt"
	"owlharbour-api/pkg/util"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	mqConn *amqp.Connection
	mqCh   *amqp.Channel
)

func CreateConnection() (*amqp.Connection, *amqp.Channel) {
	conf := mqConfig{
		Host:     util.GetEnv("RABBITMQ_HOST", ""),
		Username: util.GetEnv("RABBITMQ_USER", ""),
		Password: util.GetEnv("RABBITMQ_PASSWORD", ""),
		Vhost:    util.GetEnv("RABBITMQ_VHOST", ""),
	}

	rabbitmq := rabbitMq{mqConfig: conf}
	if mqConn == nil && mqCh == nil {
		mqConn, mqCh = rabbitmq.Connect()
	}

	fmt.Println("[*] Successfully connected to RabbitMQ.")
	return mqConn, mqCh
}
