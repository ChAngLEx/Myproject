package common

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

//NewProducer .
func NewProducer(servers string) *kafka.Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": servers})

	if err != nil {
		panic(err)
	}
	return producer
}

//Produce .
func Produce(producer *kafka.Producer, topic string, key, data []byte) error {
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          data,
		Key:            key,
	}, deliveryChan)
	if err != nil {
		return nil
	}
	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivery message to tpic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
	return nil
}

//NewConsumer .
func NewConsumer(servers string, group string, topics []string) *kafka.Consumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
		"group.id":          group,
		//"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	consumer.SubscribeTopics(topics, nil)
	return consumer
}

//Consume .
func Consume(consumer *kafka.Consumer) ([]byte, error) {
	msg, err := consumer.ReadMessage(-1)
	if err != nil {
		return nil, err
	}
	return msg.Value, nil
}
