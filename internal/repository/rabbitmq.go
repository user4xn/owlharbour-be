package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"simpel-api/internal/dto"

	"go.uber.org/zap"

	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMq interface {
	Consume(request dto.RabbitMqConsumeRequest) <-chan amqp.Delivery
	Publish(ctx context.Context, request dto.RabbitMqPublishRequest) error
	DeclareExchange(exchange dto.RabbitMQExchangeRequest)
	BindingQueue(exchangeName string, queueName string, key string)
	DeclareQueue(queueName string) amqp.Queue
}

type rabbitMq struct {
	mqConn *amqp.Connection
	mqCh   *amqp.Channel
}

func NewRabbitMqRepository(conn *amqp.Connection, ch *amqp.Channel) *rabbitMq {
	return &rabbitMq{
		mqConn: conn,
		mqCh:   ch,
	}
}

func (r *rabbitMq) Publish(ctx context.Context, request dto.RabbitMqPublishRequest) error {
	cctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	body, err := json.Marshal(request.Messages)
	if err != nil {
		fmt.Println("eFailed to marshal message", zap.String("error", err.Error()))
		return err
	}

	publish := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
		Headers:     request.Headers,
	}

	err = r.mqCh.PublishWithContext(cctx,
		request.Exchange,  // exchange
		request.QueueName, // routing key
		true,              // mandatory
		false,             // immediate
		publish,
	)

	if err != nil {
		fmt.Println("Failed to publish a message", zap.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *rabbitMq) Consume(request dto.RabbitMqConsumeRequest) <-chan amqp.Delivery {
	msgs, err := r.mqCh.Consume(
		request.QueueName,    // queue
		request.ConsumerName, // consumer
		false,                // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)

	if err != nil {
		fmt.Println("ailed to register a consumer", zap.String("error", err.Error()))
	}

	return msgs
}

func (r *rabbitMq) DeclareExchange(exchange dto.RabbitMQExchangeRequest) {
	err := r.mqCh.ExchangeDeclare(
		exchange.Name,
		exchange.Kind,
		true,
		false,
		false,
		false,
		exchange.Args,
	)

	if err != nil {
		fmt.Println(fmt.Sprintf("Failed declare exchange :`%s`", exchange.Name), zap.String("error", err.Error()))
	}

	return
}

func (r *rabbitMq) BindingQueue(exchangeName string, queueName string, key string) {
	r.DeclareQueue(queueName)

	err := r.mqCh.QueueBind(queueName, key, exchangeName, false, nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed bind queue %s to exchange `%s`", queueName, exchangeName), zap.String("error", err.Error()))
	}
}

func (r *rabbitMq) DeclareQueue(queueName string) amqp.Queue {
	q, err := r.mqCh.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Println(fmt.Sprintf("Failed declare queue: `%s`", queueName), zap.String("error", err.Error()))
	}

	return q
}
