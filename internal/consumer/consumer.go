package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/SyamSolution/notification-service/helper"
	"github.com/SyamSolution/notification-service/internal/model"
	"log"
	"os"
)

func Consumer(master sarama.Consumer, doneCh chan struct{}) {
	consumer, errors := helper.Consume(master, []string{"create-transaction", "send-pdf"})

	signals := make(chan os.Signal, 1)
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

				log.Println("consume email create transaction")

				helper.SendCreateTransactionMail(message)
			case "send-pdf":
				var message model.EmailPDFMessage
				err := json.Unmarshal(msg.Value, &message)
				if err != nil {
					fmt.Println("Error unmarshalling message", err)
				}

				log.Println("consume email send pdf")

				helper.SendEmailWithPDF(message)
			}
		case consumerError := <-errors:
			fmt.Println("Received consumer error", (consumerError).Error())
		case <-signals:
			fmt.Println("Interrupt is detected")
			doneCh <- struct{}{}
		}
	}
}
