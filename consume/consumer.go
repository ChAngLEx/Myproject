package consumer

import (
	"fmt"
	"sync"
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	
	"github.com/gomodule/redigo/redis"
	
	"github.com/olivere/elastic/v7"
	
	"gihub.com/antchfx/htmlquery"
)

var (
	indexName = "subject"
	servers   =[]string{"http://localhost:9200/"}
	client, _ = elastic.NewClient(elastic.SetURL(servers...))
)
//  定义一个Json结构体
type StdNew struct {
	Timestamp string   `json:"timestamp"`
	Source    string   `json:"source"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	URL       string   `json:"url"`
	Types     []string `json:"types"`
}

// initial ES
type LogMessage struct {
	App     string
	Topic   string
	Message string
}
var (
	esClient *elastic.Client
)

func initES(addr string) (err error) {

	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(addr))
	if err != nil{
		fmt.Println("connect es error", err)
		return
	}
	esClient = client

	return
}

func SendToES(topic string, data []byte) (err error) {
	msg := $LogMessage{}
	msg.Topic = topic
	msg.Message = string(data)

	_, err = esClient.Index().Index(topic).Type(topic).BodyJson(msg).Do()

	if err != nil{
		return
	}

	return
	
}


func PullFromKafka() {
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)

	if err != nil {
		fmt.Println("receive msg failed, err:", err)
		return
	}

	partitionList, err = consumer.Partitions("test")

	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition("test", int32(partition), sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}

		defer pc.AsyncClose()

		wg.Add(1)
		//异步从每个分区消费信息
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				fmt.Println("Partition: %d, Offset:%d, Key:%s, Value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				err = sendToES("test", mag.Value)
				if err != nil{
					logs.Warn("send to es failed, err:", err)
				}
			}
			wg.Done()
		}(pc)
		wg.Wait()
		consumer.Close()
	}
}


// 在中间件中初始化redis链接
func Redis() {
	client := redis.NewClient(&redis.Options){
		Addr:     

	}

	_, err := client.Ping().Result()

	if err != nil {
		panic(err)
	}

	RedisClient = client
}

func consumer() {
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Println("receive msg failed, err:", err)
		return
	}

	partitionList, err = consumer.Partitions("test")

	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition("test", int32(partition), sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}

		defer pc.AsyncClose()

		wg.Add(1)
		//异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				fmt.Printf("Partition: %d, Offset:%d, Key:%s, Value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(pc)
		wg.Wait()
		consumer.Close()
	}
}
