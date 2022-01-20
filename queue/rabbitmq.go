package queue

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"os"
)

type RabbitMQ struct {
	Con  *amqp.Connection
}

func NewRabbitMQ() *RabbitMQ {
	return  &RabbitMQ{}
}


func (r *RabbitMQ) Connect() (err error) {
	addr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USERNAME"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
		)

	r.Con, err = amqp.Dial(addr)
	return
}

func (r *RabbitMQ) Close() (err error) {
	err = r.Con.Close()
	return
}

func (r *RabbitMQ) Push(queueName string, body interface{}) (err error) {
	ch, err := r.Con.Channel()
	if err != nil {
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return
	}

	b, err := json.Marshal(body)

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        []byte(b),
		})

	return
}

func (r *RabbitMQ) GetMessages(queueName string) (<-chan amqp.Delivery, *amqp.Channel, error) {
	ch, err := r.Con.Channel()
	if err != nil {
		return nil, nil, err
	}


	q, err := ch.QueueDeclare(
		queueName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil,ch, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil,ch, err
	}
	return msgs,ch, nil
}
