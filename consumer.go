package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// Conectar a RabbitMQ
	conn, err := amqp.Dial("amqp://Denzel:Desz117s@18.211.110.229:5672/")
	if err != nil {
		log.Fatalf("Error conectando a RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Abrir un canal
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creando canal: %s", err)
	}
	defer ch.Close()

	// Declarar la cola (asegurarnos de que existe)
	queueName := "queue"
	q, err := ch.QueueDeclare(
		queueName,
		true,  // Durable
		false, // Auto-delete
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		log.Fatalf("Error declarando la cola: %s", err)
	}

	// Consumir mensajes de la cola
	msgs, err := ch.Consume(
		q.Name,     // Nombre de la cola
		"consumer", // Nombre del consumidor
		true,       // Auto-Acknowledge
		false,      // Exclusive
		false,      // No-local
		false,      // No-wait
		nil,        // Args
	)
	if err != nil {
		log.Fatalf("Error al consumir mensajes: %s", err)
	}

	// Canal para manejar los mensajes
	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			fmt.Printf("📥 Mensaje recibido: %s\n", msg.Body)
		}
	}()

	fmt.Println("🐰 Esperando mensajes. Presiona CTRL+C para salir...")
	<-forever // Mantiene el programa corriendo
}
