package main

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/ChAngLEx/Myproject/common"
	"github.com/garyburd/redigo/redis"
)

var globalCount = 0

//TODO isDupBloomfilter isDupSimhash
func isDupMd5(news []byte, redisclient redis.Conn) bool {
	has := md5.Sum(news)
	key := fmt.Sprintf("%x", has) //将[]byte转成16进制
	n, err := redisclient.Do("SETNX", key, "")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return n == int64(0)
}

//ConsumeToES .
func ConsumeToES() {
	consumer := common.NewConsumer("localhost:9092", "chang", []string{"test"})
	//esclient := common.NewESClient("localhost:9200")
	redisclient, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println(err.Error())
	}
	for {
		news, err := common.Consume(consumer)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		//去重 md5去重版
		if isDupMd5(news, redisclient) {
			continue
		}
		//写入ES
		//esclient.Index("news").Type("_doc").Id(key).BodyJson(news).Do(context.Background())
		globalCount++
	}
}
func main() {
	//10个并发处理
	for i := 0; i < 10; i++ {
		go ConsumeToES()
	}
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `405`)
			return
		}
		//laklahalfhajfkalgdhf
		fmt.Fprintf(w, `pong`)
	})
	http.HandleFunc("/statinfo", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `405`)
			return
		}
		//laklahalfhajfkalgdhf
		fmt.Fprintf(w, `%d`, globalCount)
	})
	http.ListenAndServe(":8081", nil)
}
