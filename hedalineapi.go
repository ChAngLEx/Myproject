package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//定义api接口数据结构和序列化json字段

type News struct {
	Timestamp uint   `json:"timestamp"`
	Source    string `json:"source`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Types     string `json:"types"`
}

//发送http请求和json 序列化并打印数据结构;

var urls = []string{
	"http://api.tianapi.com/topnews/index?key=c0ee86b42bacbd31469609cdf56dc75e",
	"http://api.tianapi.com/generalnews/index?key=c0ee86b42bacbd31469609cdf56dc75e",
	"https://3g.163.com/touch/reconstruct/article/list/BBM54PGAwangning/0-10.html",
	"https://3g.163.com/touch/reconstruct/article/list/BA10TA81wangning/0-10.html",
	"https://3g.163.com/touch/reconstruct/article/list/BA8E6OEOwangning/0-10.html",
	"https://3g.163.com/touch/reconstruct/article/list/BA8EE5GMwangning/0-10.html",
	"https://3g.163.com/touch/reconstruct/article/list/BAI67OGGwangning/0-10.html",
	"https://3g.163.com/touch/reconstruct/article/list/BA8D4A3Rwangning/0-10.html",
	"https://3g.163.com/touch/reconstruct/article/list/BAI6I0O5wangning/0-10.html",
	"https://3g.163.com/touch/reconstruct/article/list/BAI6JOD9wangning/0-10.html",
	"https://3g.163.com/touch/reconstruct/article/list/BA8F6ICNwangning/0-10.html",
	"https://3g.163.com/touch/reconstruct/article/list/BAI6RHDKwangning/0-10.html",
	"https://3g.163.com/touch/reconstruct/article/list/BA8FF5PRwangning/0-10.html",
	"https://3g.163.com/touch/reconstruct/article/list/BDC4QSV3wangning/0-10.html",
	"https://3g.163.com/touch/reconstruct/article/list/BEO4GINLwangning/0-10.html",
	"https://3g.163.com/touch/nc/api/video/recommend/Video_Recom/0-10.do?callback=videoList",
}

func main() {
	start := time.Now()
	ch := make(chan string) //创建一个字符串通道
	for _, url := range urls {
		fmt.Println("url ===", url)
		go fetch(url, ch) //go 创建协程
	}
	for range urls {
		fmt.Println(<-ch) //从通道读取信息
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	res, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	s := News{}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	secs := time.Since(start).Seconds()
	json.Unmarshal([]byte(body), &s)
	ch <- fmt.Sprintf("%.2fs %s  %s", secs, body, url) //将格式化信息写入通道
}
