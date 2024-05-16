package ship

import (
	"context"
	"encoding/json"
	"fmt"
	"owlharbour-api/internal/dto"
	"owlharbour-api/pkg/util"

	"go.uber.org/zap"
)

func (h *handler) Init() {
	exchange := dto.RabbitMQExchangeRequest{
		Name: util.GetEnv("RABBITMQ_EXCHANGE_SIMPEL_SHIP", ""),
		Kind: "direct",
	}
	h.rabbitMqRepository.DeclareExchange(exchange)
	queueName := util.GetEnv("RABBITMQ_QUEUE_SIMPEL_SHIP", "")
	h.rabbitMqRepository.DeclareQueue(queueName)
	h.rabbitMqRepository.BindingQueue(exchange.Name, queueName, "ShipRecordLog")
}

func (h *handler) WorkerRecordLog(ctx context.Context) {
	queueName := util.GetEnv("RABBITMQ_QUEUE_SIMPEL_SHIP", "")
	consumeRequest := dto.RabbitMqConsumeRequest{
		QueueName:    queueName,
		ConsumerName: "execute-ship-record-log",
	}

	msgs := h.rabbitMqRepository.Consume(consumeRequest)

	forever := make(chan struct{})
	go func() {
		defer close(forever)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Context cancelled, exiting WorkerRecordLog")
				return
			case m, ok := <-msgs:
				if !ok {
					return // Channel closed, exit goroutine
				}
				var data dto.ShipRecordRequest
				err := json.Unmarshal(m.Body, &data)
				if err != nil {
					fmt.Println("Failed to unmarshal message", zap.String("error", err.Error()))
					m.Ack(true)
					continue
				}

				err = h.service.RecordLocationShip(ctx, data)

				if err != nil {
					fmt.Println("Error processing ship log", zap.String("device id", data.DeviceID), zap.String("error :", err.Error()))
					m.Ack(false)
					continue
				}

				m.Ack(true)
			}
		}
	}()
	fmt.Println("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
