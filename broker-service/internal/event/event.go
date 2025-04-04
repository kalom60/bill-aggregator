package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"refresh_bills", // name
		"topic",         // type
		true,            // durable?
		false,           // auto-deleted?
		false,           // internal?
		false,           // no-wait?
		nil,             // arguments?
	)
}
