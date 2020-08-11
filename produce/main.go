package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ChAngLEx/Myproject/common"
)

var sources = map[string]common.NewsSource{
	"TopNews":     common.NewsSource{URL: "http://api.tianapi.com/topnews/index?key=c0ee86b42bacbd31469609cdf56dc75e", Parse: TopNews2StdNew},
	"GeneralNews": common.NewsSource{URL: "http://api.tianapi.com/generalnews/index?key=c0ee86b42bacbd31469609cdf56dc75e", Parse: GeneralNews2StdNew},
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
func TopNews2StdNew(topnews []byte) ([]common.StdNew, error) {
	var t TopNews
	err := json.Unmarshal(topnews, &t)
	if err != nil {
		return nil, err
	}
	ret := []common.StdNew{}
	for i := range t.Newslist {
		var item common.StdNew
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
func GeneralNews2StdNew(generalnews []byte) ([]common.StdNew, error) {
	var g GeneralNews
	err := json.Unmarshal(generalnews, &g)
	if err != nil {
		return nil, err
	}
	ret := []common.StdNew{}
	for i := range g.Newslist {
		var item common.StdNew
		item.Source = "general"
		item.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		item.Title = g.Newslist[i].Title
		item.Body = g.Newslist[i].Description
		item.URL = g.Newslist[i].URL
		ret = append(ret, item)
	}
	return ret, nil
}

var dataCount = 0

func PullAllAPI() []common.StdNew {
	start := time.Now()
	wg := sync.WaitGroup{}
	for srcname, srcinfo := range sources {
		wg.Add(1)
		go func(srcname string, srcinfo common.NewsSource) {
			log.Printf("name:%s\turl:%s\n", srcname, srcinfo.URL)
			results, _ := fetch(srcinfo)
			for i := range results {
				fmt.Println("%v", results[i])
				ProduceToKafka(results[i])
				dataCount++
				log.Printf("%v", results[i])
			}
			wg.Done()
		}(srcname, srcinfo)
	}
	wg.Wait()
	log.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
	return nil

}

//ProduceToKafka 将StdNew写入kafka
func ProduceToKafka(news common.StdNew) {
	producer := common.NewProducer("localhost:9092")

	for {
		value, _ := json.Marshal(news)
		err := common.Produce(producer, "test", nil, []byte(value))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("send success")
		}

	}

}

func main() {

	for i := 0; i < 10; i++ {
		go PullAllAPI()
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `405`)
			return
		}
		fmt.Fprintf(w, `pong`)
	})
	http.ListenAndServe(":80", nil)

	http.HandleFunc("/statinfo", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `405`)
			return
		}
		//laklahalfhajfkalgdhf
		fmt.Fprintf(w, `%d`, dataCount)
	})
	http.ListenAndServe(":8081", nil)

}

func fetch(srcinfo common.NewsSource) ([]common.StdNew, error) {
	start := time.Now()
	res, err := http.Get(srcinfo.URL)
	if err != nil {

	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body) //TODO 把这个数据存到文件里，因为只有1000次免费调用
	log.Printf("%s cost:%f\n", srcinfo.URL, time.Since(start).Seconds())
	return srcinfo.Parse(body)
}
