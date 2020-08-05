# Myproject
Go/API/Data Processing

Project Develop and Design 

TECH: GO<Docker<Kafka<Json<Linux<Git<ElasticSearch<Redis

Step1: Data Access & Unified data Structure & Push data into Kafka

Call APIs from different sources to filter, de-duplicate, aggregate, and select article data, as well as article retrieval, recommendation.
Achieve data access from different sources, and normalizing the data into a unified format, storing it in the database and pushing it downstream

Sources: 
1. news API:         https://www.juhe.cn/docs/api/id/235
2. 163News API:       https://www.jianshu.com/p/c54e25349b77
3. TOPNEWS API:      https://www.tianapi.com/apiview/99
4. GENERNALNEWS API: https://www.tianapi.com/apiview/87

Unified Data Structure:

type StdNew struct {
	Timestamp string   `json:"timestamp"`
	Source    string   `json:"source"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	URL       string   `json:"url"`
	Types     []string `json:"types"`
}



Step2: Pull data From Kafka & Deduplicate & Send to Elasticsearch (Filter Data) 

Filtering(such as profanity), de-duplication and marking (keyword extraction, category labeling) of normalized data. Push downstream in batches.
Used Kafka to produce and consume data. Used Redis to filter data and write data into Elasticsearch, 





Step3: Conncect to Web service and Display

Responsible for implementing back-end interfaces such as article data retrieval, recommendation, and public timeline pull, and design and implement Web pages to display various functions.

accomplish the connection between web service. 


