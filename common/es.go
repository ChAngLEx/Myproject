package common

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/olivere/elastic"
)

//NewESClient .
func NewESClient(host string) *elastic.Client {
	client, err := elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetGzip(true),
		elastic.SetSniff(false),
		elastic.SetHttpClient(&http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        200,
				IdleConnTimeout:     30 * time.Second,
				MaxConnsPerHost:     100,
				MaxIdleConnsPerHost: 100,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
			},
		}),
	)
	if err != nil {
		panic(err)
	}
	return client
}
func after(executionID int64, requests []elastic.BulkableRequest, response *elastic.BulkResponse, err error) {
	if err != nil {
		log.Printf("bulk commit failed, err: %v\n", err)
	}
	// do what ever you want in case bulk commit success
	log.Printf("commit successfully, len(requests)=%d\n", len(requests))
}

//NewBulkProcessor .
func NewBulkProcessor(client *elastic.Client) *elastic.BulkProcessor {
	bulkProcessor, err := elastic.NewBulkProcessorService(client).
		Workers(5).
		BulkActions(1000).
		FlushInterval(1 * time.Second).
		After(after).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	return bulkProcessor
}
