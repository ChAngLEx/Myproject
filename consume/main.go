package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"gopkg.in/jcmturner/gokrb5.v7/config"

	"github.com/olivere/elastic/v7"
)

// log initial
type LogConfig struct {
	kafkaAddr string
	ESAddr    string
	LogPath   string
	LogLevel  string
	Topic     string
}

var (
	logConfig *LogConfig
)

func initConfig(conftype string, filename string) (err error) {
	conf, err := config.NewConfig(conftype, filename)
	if err != nil {
		fmt.Println("new config faild,err:", err)
		return
	}

	logConfig = &LogConfig{}
	logConfig.LogLevel = conf.String("logs::log_level")
	if len(logConfig.LogLevel) == 0 {
		logConfig.LogLevel = "debug"
	}

	logConfig.LogPath = conf.String("logs::log_path")
	if len(logConfig.LogPath) == 0 {
		logConfig.LogPath = "./logs"
	}

	logConfig.kafkaAddr = conf.String("kafka::server_addr")
	if len(logConfig.kafkaAddr) == 0 {
		err = fmt.Errorf("invalid kafka addr err")
		return
	}

	logConfig.ESAddr = conf.String("ES::addr")
	if len(logConfig.ESAddr) == 0 {
		err = fmt.Errorf("invalid ES addr err")
		return
	}

	logConfig.Topic = conf.String("kafka::topic")
	if len(logConfig.Topic) == 0 {
		err = fmt.Errorf("invalid topic addr err")
		return
	}
	return
}

func convertLogLevel(level string) int {

	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}

	return logs.LevelDebug
}

func initLogger(logpath string, logLevel string) (err error) {
	config := make(map[string]interface{})
	config["filename"] = logpath
	config["level"] = convertLogLevel(logLevel)

	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("marshal failed,err:", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}

// initial kafka
type KafkaClient struct {
	client sarama.Consumer
	addr   string
	topic  string
	wg     sync.WaitGroup
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

type LogMessage struct {
	App     string
	Topic   string
	Message string
}

var (
	esClient *elastic.Client
)

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

func sendToES(topic string, data []byte) (err error) {
	msg := &LogMessage{}
	msg.Topic = topic
	msg.Message = string(data)

	_, err = esClient.Index().
		Index(topic).
		Type(topic).
		BodyJson(msg).
		Do()
	if err != nil {

		return
	}

	return
}
