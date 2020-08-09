 Data Processing Project (Go/Golang)
 ===



A Go Application combined API service, Apache Kafka, and Elasticsearch,  which handles news synchronization and news retrieval. 

### Technologies ###
Go/Golang || HTTP Service || Apche Kafka || Redis || Elasticsearch || Linux





### Process Diagram
<div align=center><img width="625" height="475" src="https://github.com/ChAngLEx/Myproject/blob/master/image/process%20Diagram.jpg"/></div>
																       
    








### Step1 ： Data Access and Unified ###

<img width="325" height= "225" src = "https://github.com/ChAngLEx/Myproject/blob/master/image/process%20diagram%20part1.jpg"/>

1 Connect different open news sources with API service.<br/>
2 Normalizing the data into a unified format. <br/>
3 Feed data from the Apache Kafka and pushing it downstream.<br/>
```
Open Sources Example: 
1. news API:          https://www.juhe.cn/docs/api/id/235
2. 163News API:       https://www.jianshu.com/p/c54e25349b77
3. TOPNEWS API:       https://www.tianapi.com/apiview/99
4. GENERNALNEWS API:  https://www.tianapi.com/apiview/87

```

```

Unified Data Structure:

type StdNew struct {
	Timestamp string   `json:"timestamp"`
	Source    string   `json:"source"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	URL       string   `json:"url"`
	Types     []string `json:"types"`
}
```
#### References ###
Kafka repository: https://github.com/confluentinc/confluent-kafka-go<br/>
API Crawling:     https://www.tianapi.com/apiview/87<br/>





### Step2: Data optimization and Organization ###

<img width="275" height= "125" src = "https://github.com/ChAngLEx/Myproject/blob/master/image/process%20diagram%20part2.jpg"/>

1 Consume data from Apache Kafka.<br/>
2 Optimize data, using Redis filtering news to accomplish data deduplication.<br/>
3 Create the data index and write into Elasticsearch.<br/>

```
Open Source:
Content Moderation API: 
1. Tencent could:                https://cloud.tencent.com/product/tms
2. Baidu could:                  https://ai.baidu.com/tech/textcensoring
Others:
1. jieba:                        https://github.com/fxsjy/jieba 
2. Sensitive words:              https://github.com/jkiss/sensitive-words
3. Trie tree with AC automaton:  https://github.com/ChisBread/TrieNAhoCorasick
4. Tencent Cloud NLP:            https://cloud.tencent.com/product/nlp
5. simhash:                      https://blog.csdn.net/lengye7/article/details/79789206


```

#### References ###
Redis repository:         https://godoc.org/github.com/gomodule/redigo/redis#pkg-index<br/>
Elasticsearch repository: https://github.com/olivere/elastic<br/>
Library:                  https://godoc.org/ <br/>





### Step3: Connect to Web service and Display ###
<img width="275" height = "175" src="https://github.com/ChAngLEx/Myproject/blob/master/image/process%20diagram%20part3.jpg"/>

1 Implement back-end interfaces for news retrieval, news recommendation, and news timeline count.<br/>
2 Front-end developer is responsible for the UI design.<br/>

 


