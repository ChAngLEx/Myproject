### Data Processing Project (Go/Golang) ###
A Go Application combined API service, Apache Kafka, and Elasticsearch,  which handles news synchronization and news retrieval. 





### Process Diagram

    








### Step1 ： Data Access and Unified ###

1 Connect different open news sources with API service.
2 Normalizing the data into a unified format. 
3 Feed data from the Apache Kafka and pushing it downstream.

```
Open Sources Example: 
1. news API:          https://www.juhe.cn/docs/api/id/235
2. 163News API:       https://www.jianshu.com/p/c54e25349b77
3. TOPNEWS API:       https://www.tianapi.com/apiview/99
4. GENERNALNEWS API:  https://www.tianapi.com/apiview/87

```

'''

Unified Data Structure:

type StdNew struct {
	Timestamp string   `json:"timestamp"`
	Source    string   `json:"source"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	URL       string   `json:"url"`
	Types     []string `json:"types"`
}

'''



### Step2: Data optimization and Organization ###

1 Consume data from Apache Kafka
2 Optimize data, using Redis filtering news to accomplish data deduplication.
3 Create the data index and write into Elasticsearch






### Step3: Connect to Web service and Display ###

1 Implement back-end interfaces for news retrieval, news recommendation, and news timeline count.
2 Front-end developer is responsible for the UI design.

 


