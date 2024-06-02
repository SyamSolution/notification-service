package helper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)


func processMessage(message types.Message) {
    fmt.Printf("Message ID: %s\n", *message.MessageId)
    fmt.Printf("Message Body: %s\n", *message.Body)
}

func StartConsumer() {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }

    client := sqs.NewFromConfig(cfg)

    queueURL := "https://sqs.ap-southeast-1.amazonaws.com/735185491450/syam-message"

    for {
        result, err := client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
            QueueUrl:            &queueURL,
            MaxNumberOfMessages: 10,
            WaitTimeSeconds:     10,
        })
        if err != nil {
            log.Printf("failed to receive messages, %v", err)
            continue
        }

        if len(result.Messages) == 0 {
            fmt.Println("No messages received")
            continue
        }

        for _, message := range result.Messages {
            processMessage(message)

            _, err := client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
                QueueUrl:      &queueURL,
                ReceiptHandle: message.ReceiptHandle,
            })
            if err != nil {
                log.Printf("failed to delete message, %v", err)
            }
        }

        time.Sleep(1 * time.Second)
    }
}
