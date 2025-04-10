package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

const (
	rabbitMQURL = "amqp://Denzel:Desz117s@18.211.110.229:5672/"
	queueName   = "queue"
	apiURL      = "http://3.86.136.69:8001/pedidos/log"
)

func main() {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Error conectando a RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creando canal: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error declarando la cola: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"consumer",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error al consumir mensajes: %s", err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			var pedido map[string]interface{}

			err := json.Unmarshal(msg.Body, &pedido)
			if err != nil {
				log.Printf("❌ Error deserializando mensaje: %s\n", err)
				continue
			}

			prettyJSON, _ := json.MarshalIndent(pedido, "", "  ")
			fmt.Println("📩 Pedido recibido desde RabbitMQ:")
			fmt.Println(string(prettyJSON))
			fmt.Println("----------------------------")

			sendToAPI(pedido)
		}
	}()

	fmt.Println("🐇 Esperando mensajes de la cola. Presiona CTRL+C para salir...")
	<-forever
}

func sendToAPI(data map[string]interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("❌ Error serializando JSON para la API: %s\n", err)
		return
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ Error enviando datos a la API: %s\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Println("✅ Pedido enviado y recibido correctamente en la API2")
	} else {
		log.Printf("⚠️ Respuesta inesperada de la API2. Código: %d\n", resp.StatusCode)
	}
}
