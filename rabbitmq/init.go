package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type rbMQ struct {
	messages chan []byte
	queue    string
}

func newRabbitMQ(queue string) *rbMQ {
	return &rbMQ{make(chan []byte), queue}
}

func executor(algo func(ch *amqp.Channel, rmq *rbMQ), rmq *rbMQ) {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		panic(1)
	}
	ch, _ := conn.Channel()
	algo(ch, rmq)

	defer conn.Close()
}

func producer(ch *amqp.Channel, r *rbMQ) {
	q, err := ch.QueueDeclare(
		r.queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(<-r.messages),
		})
	failOnError(err, "Failed to publish a message")

}

func consumer(ch *amqp.Channel, r *rbMQ) {
	q, err := ch.QueueDeclare(
		r.queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			r.messages <- d.Body
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func FromQueque(mq *rbMQ) {
	executor(consumer, mq)
}

func toQueque(mq *rbMQ) {
	executor(producer, mq)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
