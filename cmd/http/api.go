package main

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/SyamSolution/notification-service/config"
	"github.com/SyamSolution/notification-service/helper"
	"github.com/SyamSolution/notification-service/model"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	baseDep := config.NewBaseDep()
	loadEnv(baseDep.Logger)

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	brokers := []string{os.Getenv("KAFKA_BROKER")}
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Panicf("Error creating consumer: %s", err)
	}
	defer func() {
		if err := master.Close(); err != nil {
			log.Panicf("Error closing consumer: %s", err)
		}
	}()

	log.Println("Connected to Kafka broker")

	consumer, errors := helper.Consume(master, []string{"create-transaction", "completed-transaction"})

	signals := make(chan os.Signal, 1)

	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-consumer:
				switch msg.Topic {
				case "create-transaction":
					var message model.DataMessage
					err := json.Unmarshal(msg.Value, &message)
					if err != nil {
						fmt.Println("Error unmarshalling message", err)
					}

					helper.SendCreateTransactionMail(message)
				case "completed-transaction":
					var message model.CompleteTransactionMessage
					err := json.Unmarshal(msg.Value, &message)
					if err != nil {
						fmt.Println("Error unmarshalling message", err)
					}

					log.Println(message)

					helper.SendCompletedTransactionMail(message)
				}
			case consumerError := <-errors:
				fmt.Println("Received consumer error", (consumerError).Error())
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
}

func loadEnv(logger config.Logger) {
	_, err := os.Stat(".env")
	if err == nil {
		err = godotenv.Load()
		if err != nil {
			logger.Error("no .env files provided")
		}
	}
}
