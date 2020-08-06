package main

import (
	"fmt"
	"strings"

	"github.com/Shopify/sarama"

	"github.com/olivere/elastic/v7"
)

type Service struct {
	AccessLogProducer sarama.AsyncProducer
	Consumer          sarama.Consumer
	Client            *elastic.Client
}

var (
	kafkaClient *KafkaClient
)

func initKafKa(addr string, topic string) (err error) {

	kafkaClient = &KafkaClient{}

	consumer, err := sarama.NewConsumer(strings.Split(addr, ","), nil)
	if err != nil {
		logs.Error("Failed to strat consumer :", err)
		return
	}

	kafkaClient.client = consumer
	kafkaClient.addr = addr
	kafkaClient.topic = topic
	return

}

var (
	esClient *elastic.Client
)

func PullFromKafka() (err error) {
	partitionList, err := kafkaClient.client.Partitions(kafkaClient.topic)
	if err != nil {
		logs.Error("ini failed ,err:%v", err)
		fmt.Printf("ini failed ,err:%v", err)
		return
	}

	for partition := range partitionList {
		fmt.Println("for in")
		pc, errRet := kafkaClient.client.ConsumePartition(kafkaClient.topic, int32(partition), sarama.OffsetNewest)
		if errRet != nil {
			err = errRet
			logs.Error("Failed to start consumer for partition %d: %s\n", partition, err)
			fmt.Printf("Failed to start consumer for partition %d: %s\n", partition, err)
			return
		}
		defer pc.AsyncClose()

		kafkaClient.wg.Add(1)
		go func(pc sarama.PartitionConsumer) {

			for msg := range pc.Messages() {

				logs.Debug("Partition:%d,Offset:%d,key:%s,value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				//fmt.Println()
				err = sendToES(kafkaClient.topic, msg.Value)
				if err != nil {
					logs.Warn("send to es failed,err:%v", err)
				}
			}
			kafkaClient.wg.Done()
		}(pc)

	}

	kafkaClient.wg.Wait()
	return
}

// initial ES
func initES(addr string) (err error) {

	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(addr))
	if err != nil {
		fmt.Println("connect es error", err)
		return
	}
	esClient = client

	return
}

func sendToES(body, []byte) error {
	ctx := context.Backgroud()

	//topic structure need be defined?
	put, err = esClient.Index().
		Index(topic).
		Type(topic).
		BodyJson(body).
		Do(ctx)
	if err != nil {

		panic(err)
	}

	fmt.Println("Indexed log %s to index %s, type %s\n", put.Id, put.Index, put.Type)

	// index need a structure...
	_, err = esClient.Flush().Index("").Do(ctx)
	if err != nil {
		panic(err)

	}
	return nil
}
