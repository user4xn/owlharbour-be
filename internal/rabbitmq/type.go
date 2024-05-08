package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	mqConfig struct {
		Username string
		Password string
		Vhost    string
		Host     string
	}

	rabbitMq struct {
		mqConfig
	}
)

func (conf rabbitMq) Connect() (*amqp.Connection, *amqp.Channel) {
	var err error
	connStr := fmt.Sprintf("amqp://%s:%s@%s/%s", conf.Username, conf.Password, conf.Host, conf.Vhost)
	conn, err := amqp.Dial(connStr)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return conn, ch
}
