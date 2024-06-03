package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/SyamSolution/notification-service/helper"
	"github.com/SyamSolution/notification-service/internal/model"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const workerCount = 5

func processCreateTransactionMessage(message types.Message) {
    fmt.Printf("CreateTransaction Queue - Message ID: %s\n", *message.MessageId)
    fmt.Printf("CreateTransaction Queue - Message Body: %s\n", *message.Body)
}

func processSendEmailPdfMessage(message types.Message) {
    fmt.Printf("SendEmailPdf Queue - Message ID: %s\n", *message.MessageId)
    fmt.Printf("SendEmailPdf Queue - Message Body: %s\n", *message.Body)
}

func processDeadLetterMessage(message types.Message) {
    fmt.Printf("Dead Letter Queue - Message ID: %s\n", *message.MessageId)
    fmt.Printf("Dead Letter Queue - Message Body: %s\n", *message.Body)
}

func workerDeadLetter(client *sqs.Client, queueURL string, wg *sync.WaitGroup, dqlType string) {
    defer wg.Done()
    for {
        result, err := client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
            QueueUrl:            &queueURL,
            MaxNumberOfMessages: 10,
            WaitTimeSeconds:     10,
        })
        if err != nil {
            log.Printf("failed to receive messages from Dead Letter queue, %v", err)
            continue
        }

        if len(result.Messages) == 0 {
            time.Sleep(1 * time.Second)
            continue
        }

        for _, message := range result.Messages {
            if dqlType == "transaction" {
                var msg model.DataMessage
			    err := json.Unmarshal([]byte(*message.Body), &msg)
				if err != nil {
					fmt.Println("Error unmarshalling message", err)
				}else{
                    log.Println("consume email create transaction")
                    helper.SendCreateTransactionMail(msg)
                }
            } else if dqlType == "pdf" {
                var msg model.EmailPDFMessage
                err := json.Unmarshal([]byte(*message.Body), &msg)
                if err != nil {
                    fmt.Println("Error unmarshalling message", err)
                }else{
                    log.Println("consume email send pdf")
                    helper.SendEmailWithPDF(msg)
                }
            }
            processDeadLetterMessage(message)

            _, err := client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
                QueueUrl:      &queueURL,
                ReceiptHandle: message.ReceiptHandle,
            })
            if err != nil {
                log.Printf("failed to delete message from Dead Letter queue, %v", err)
            }
        }
    }
}

func workerCreateTransaction(client *sqs.Client, queueURL string, wg *sync.WaitGroup) {
    defer wg.Done()
    for {
        result, err := client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
            QueueUrl:            &queueURL,
            MaxNumberOfMessages: 10,
            WaitTimeSeconds:     10,
        })
        if err != nil {
            log.Printf("failed to receive messages from CreateTransaction queue, %v", err)
            continue
        }

        if len(result.Messages) == 0 {
            time.Sleep(1 * time.Second)
            continue
        }

        for _, message := range result.Messages {
		    var msg model.DataMessage
			    err := json.Unmarshal([]byte(*message.Body), &msg)
				if err != nil {
					fmt.Println("Error unmarshalling message", err)
				}else {
                    log.Println("consume email create transaction")
                    helper.SendCreateTransactionMail(msg)
                }

				
            processCreateTransactionMessage(message)

            _, err = client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
                QueueUrl:      &queueURL,
                ReceiptHandle: message.ReceiptHandle,
            })
            if err != nil {
                log.Printf("failed to delete message from CreateTransaction queue, %v", err)
            }
        }
    }
}

func workerSendEmailPdf(client *sqs.Client, queueURL string, wg *sync.WaitGroup) {
    defer wg.Done()
    for {
        result, err := client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
            QueueUrl:            &queueURL,
            MaxNumberOfMessages: 10,
            WaitTimeSeconds:     10,
        })
        if err != nil {
            log.Printf("failed to receive messages from SendEmailPdf queue, %v", err)
            continue
        }

        if len(result.Messages) == 0 {
            time.Sleep(1 * time.Second)
            continue
        }

        for _, message := range result.Messages {
            var msg model.EmailPDFMessage
            err := json.Unmarshal([]byte(*message.Body), &msg)
            if err != nil {
                fmt.Println("Error unmarshalling message", err)
            }else{
                log.Println("consume email send pdf")
                helper.SendEmailWithPDF(msg)
            }
            processSendEmailPdfMessage(message)

            _, err = client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
                QueueUrl:      &queueURL,
                ReceiptHandle: message.ReceiptHandle,
            })
            if err != nil {
                log.Printf("failed to delete message from SendEmailPdf queue, %v", err)
            }
        }
    }
}

func StartConsumer() {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }

    client := sqs.NewFromConfig(cfg)

    var wg sync.WaitGroup

    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go workerCreateTransaction(client, os.Getenv("SQS_TRANSACTION_URL"), &wg)
    }

    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go workerSendEmailPdf(client, os.Getenv("SQS_MAIL_URL"), &wg)
    }

    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go workerDeadLetter(client, os.Getenv("SQS_TRANSACTION_DLQ_URL"), &wg, "transaction")
    }

    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go workerDeadLetter(client, os.Getenv("SQS_MAIL_DLQ_URL"), &wg, "pdf")
    }

    wg.Wait()
}
