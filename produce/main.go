package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

//StdNew 标准格式
type StdNew struct {
	Timestamp string   `json:"timestamp"`
	Source    string   `json:"source"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	URL       string   `json:"url"`
	Types     []string `json:"types"`
}

//NewsSource URL:API Parse:API结果转StdNews的函数
type NewsSource struct {
	URL   string
	Parse func([]byte) ([]StdNew, error) //函数可以作为变量
}

var sources = map[string]NewsSource{
	"TopNews":     NewsSource{URL: "http://api.tianapi.com/topnews/index?key=c0ee86b42bacbd31469609cdf56dc75e", Parse: TopNews2StdNew},
	"GeneralNews": NewsSource{URL: "http://api.tianapi.com/generalnews/index?key=c0ee86b42bacbd31469609cdf56dc75e", Parse: GeneralNews2StdNew},
}

//TopNews TopNews返回结果
type TopNews struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	Newslist []struct {
		Ctime       string `json:"ctime"`
		Title       string `json:"title"`
		Description string `json:"description"`
		PicURL      string `json:"picUrl"`
		URL         string `json:"url"`
		Source      string `json:"source"`
	} `json:"newslist"`
}

//TopNews2StdNew topnews转成[]StdNew
func TopNews2StdNew(topnews []byte) ([]StdNew, error) {
	var t TopNews
	err := json.Unmarshal(topnews, &t)
	if err != nil {
		return nil, err
	}
	ret := []StdNew{}
	for i := range t.Newslist {
		var item StdNew
		item.Source = "top"
		item.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		item.Title = t.Newslist[i].Title
		item.Body = t.Newslist[i].Description
		item.URL = t.Newslist[i].URL
		ret = append(ret, item)
	}
	return ret, nil
}

type GeneralNews struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	Newslist []struct {
		Ctime       string `json:"ctime"`
		Title       string `json:"title"`
		Description string `json:"description"`
		PicURL      string `json:"picUrl"`
		URL         string `json:"url"`
		Source      string `json:"source"`
	} `json:"newslist"`
}

//GeneralNews2StdNew generalnews转成[]StdNew
func GeneralNews2StdNew(generalnews []byte) ([]StdNew, error) {
	//TODO 和上面差不多，印象中GeneralNews和TopNews是一个格式
	var g GeneralNews
	err := json.Unmarshal(generalnews, &g)
	if err != nil {
		return nil, err
	}
	ret := []StdNew{}
	for i := range g.Newslist {
		var item StdNew
		item.Source = "general"
		item.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		item.Title = g.Newslist[i].Title
		item.Body = g.Newslist[i].Description
		item.URL = g.Newslist[i].URL
		ret = append(ret, item)
	}
	return ret, nil
}

//CallAPI 调用API, 将转换好的StdNew存入chan中
func CallAPI(ch chan StdNew) {
	start := time.Now()
	//进行一轮并发API调用
	wg := sync.WaitGroup{}
	for srcname, srcinfo := range sources {
		wg.Add(1)
		go func(srcname string, srcinfo NewsSource) {
			log.Printf("name:%s\turl:%s\n", srcname, srcinfo.URL)
			results, _ := fetch(srcinfo)
			for i := range results {
				//TODO kafka produce
				PushToKafka(results[i])
				log.Printf("%v", results[i])
			}
			wg.Done()
		}(srcname, srcinfo)
	}
	wg.Wait()
	log.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

//PushToKafka 将StdNew写入kafka
func PushToKafka(new StdNew) {
	config := sarama.NewConfig() // 配置
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second

	//  构造消息
	//msg := &sarama.ProducerMessage{}
	//msg.Topic = "test"
	//msg.Value = sarama.StringEncoder("this is a test")

	// connect kafka
	p, err := sarama.NewSyncProducer([]string{"127.0.0.1: 9092"}, config)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return
	}
	defer p.Close()

	data, err := json.Marshal(new)

	for i := 0; i < 10; i++ {
		msg := &sarama.ProducerMessage{
			Topic: "test",
			Value: sarama.ByteEncoder(data),
		}

		if part, offset, err := p.SendMessage(msg); err != nil {
			fmt.Fprintf(os.Stdout, "Send Failed, Error: %v \n", err)
		} else {
			fmt.Fprintf(os.Stdout, "Send Success, partition=%d, offset=%d \n", part, offset)
		}
	}
}

func main() {
	ch := make(chan StdNew, 1024) //带buffer的channel不阻塞生产和消费
	CallAPI(ch)                   //生产原始数据
	//TODO 改成循环CallAPI

	/*for i := 0; i < 10; i++ {
		CallAPI(ch)
	}
	for i := 0; i < 10; i++ {
		<-ch
	}*/

	//TODO 实现从ch中读数据，并调用PushToKafka

}

func fetch(srcinfo NewsSource) ([]StdNew, error) {
	start := time.Now()
	res, err := http.Get(srcinfo.URL)
	if err != nil {

	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body) //TODO 把这个数据存到文件里，因为只有1000次免费调用
	log.Printf("%s cost:%f\n", srcinfo.URL, time.Since(start).Seconds())
	return srcinfo.Parse(body)
}
