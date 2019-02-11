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
			config.GetConfig().AMQP.Username,
			config.GetConfig().AMQP.Password,
			config.GetConfig().AMQP.Addr,
		),
	)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %s", err.Error())
		return err
	}
	log.Info("successfully connect to RabbitMQ")

	// open a consumer channel
	ConsumerChannel, err = ConsumerConnection.Channel()
	if err != nil {
		log.Fatalf("failed to open a consumer channel: %s", err.Error())
		return err
	}
	log.Info("successfully open consumer")

	// declare a consumer queue
	ConsumerQueue, err = ConsumerChannel.QueueDeclare(
		config.GetConfig().AMQP.QueueName,
		config.GetConfig().AMQP.QueueDurable,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare a consumer queue: %s", err.Error())
		return err
	}
	log.Infof("successfully declare consumer queue with name %s", ConsumerQueue.Name)

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
	log.Info("successfully register consumer")

	return nil
}
