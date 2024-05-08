package dto

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	RabbitMqUpdateOrderStatusRequest struct {
		OrderNo string `json:"order_no"`
		Status  string `json:"status"`
		Cnote   string `json:"cnote"`
		UserId  int    `json:"user_id"`
	}

	RabbitMqPublishRequest struct {
		Exchange  string
		QueueName string
		Headers   amqp.Table
		Messages  interface{}
	}

	RabbitMqConsumeRequest struct {
		QueueName    string
		ConsumerName string
	}

	RabbitMQExchangeRequest struct {
		Name string
		Kind string
		Args amqp.Table
	}
)
