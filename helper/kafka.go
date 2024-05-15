package helper

import (
	"github.com/IBM/sarama"
	"log"
)

func Consume(consumer sarama.Consumer, topics []string) (<-chan *sarama.ConsumerMessage, <-chan *sarama.ConsumerError) {
	messages := make(chan *sarama.ConsumerMessage, 256)
	errors := make(chan *sarama.ConsumerError, 256)

	for _, topic := range topics {
		partitionList, err := consumer.Partitions(topic)
		if err != nil {
			log.Printf("Error retrieving partition list for topic %s: %s\n", topic, err)
			continue
		}

		for _, partition := range partitionList {
			pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				log.Printf("Error consuming partition %d of topic %s: %s\n", partition, topic, err)
				continue
			}

			go func(pc sarama.PartitionConsumer) {
				for message := range pc.Messages() {
					messages <- message
				}
			}(pc)

			go func(pc sarama.PartitionConsumer) {
				for err := range pc.Errors() {
					errors <- err
				}
			}(pc)
		}
	}

	return messages, errors
}
