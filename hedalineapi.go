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
	Code     int               `json:"code"`
	Message  string            `json:"msg"`
	Newslist map[string]string `json:"newslist"`
}

type Newslist struct {
	Timestamp uint   `json:"ctime"`
	Source    string `json:"source`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Types     string `json:"types"`
}

//发送http请求和json 序列化并打印数据结构;

var urls = []string{
	"http://api.tianapi.com/topnews/index?key=c0ee86b42bacbd31469609cdf56dc75e",
	"http://api.tianapi.com/generalnews/index?key=c0ee86b42bacbd31469609cdf56dc75e",
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

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	secs := time.Since(start).Seconds()
	fmt.Println(<-s)
	ch <- fmt.Sprintf("%.2fs %s  %s", secs, body, url) //将格式化信息写入通道
}

func InterfaceMethod() {
	b := []byte(`{`)

	var f interface{}
	_ = json.Unmarshal([]byte(b), &f)
	//err := json.Unmarshal(b, &f)

	m := f.(map[string]interface{})

	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case int:
			fmt.Println(k, "is int", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			for i, u := range vv {
				fmt.Println(i, u)
			}
		case map[string]interface{}:
			fmt.Println(k, "is an map[string]string:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't handle")
		}
	}
}
