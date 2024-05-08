package model

import amqp "github.com/rabbitmq/amqp091-go"

type (
	RabbitMQExchange struct {
		Name string
		Kind string
		Args amqp.Table
	}
)
