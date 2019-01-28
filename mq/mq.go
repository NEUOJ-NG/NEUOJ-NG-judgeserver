package mq

import (
	"fmt"
	"github.com/NEUOJ-NG/NEUOJ-NG-judgeserver/config"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	ConsumerConnection *amqp.Connection
	ConsumerChannel    *amqp.Channel
	ConsumerQueue      amqp.Queue
	ConsumerMessages   <-chan amqp.Delivery
)

// initialize consumer message queue
// remember to close connections and channels
// after use outside the function
func InitConsumerMQ() error {
	var err error

	// connect
	log.Info("connecting to RabbitMQ")
	ConsumerConnection, err = amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s/",
			config.GetConfig().AMQPConfig.Username,
			config.GetConfig().AMQPConfig.Password,
			config.GetConfig().AMQPConfig.Addr,
		),
	)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %s", err.Error())
		return err
	}
	log.Info("connect to RabbitMQ success")

	// open a consumer channel
	ConsumerChannel, err = ConsumerConnection.Channel()
	if err != nil {
		log.Fatalf("failed to open a consumer channel: %s", err.Error())
		return err
	}
	log.Info("open consumer success")

	// declare a consumer queue
	ConsumerQueue, err = ConsumerChannel.QueueDeclare(
		config.GetConfig().AMQPConfig.QueueName,
		config.GetConfig().AMQPConfig.QueueDurable,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare a consumer queue: %s", err.Error())
		return err
	}
	log.Infof("declare consumer queue with name %s success", ConsumerQueue.Name)

	// register consumer
	ConsumerMessages, err = ConsumerChannel.Consume(
		ConsumerQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to register a consumer: %s", err.Error())
		return err
	}
	log.Info("register consumer success")

	return nil
}
