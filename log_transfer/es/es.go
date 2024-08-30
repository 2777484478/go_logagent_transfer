package es

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io"
	"strings"
	"time"
)

type LogData struct {
	Data  string `json:"data"`
	Topic string `json:"topic"`
}

var (
	client *elasticsearch.Client
	ch     chan *LogData
)

func Init(address string) (err error) {
	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	fmt.Println(address)
	cfg := elasticsearch.Config{
		Addresses: []string{
			address,
		},
	}
	client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Printf("Error creating the Elasticsearch client: %s", err)
	}
	fmt.Println("connect to ES ")
	ch = make(chan *LogData, 100000)
	go SendTOES1()
	return
}

func SendTOES(index string, data string) {

	fmt.Println(index, data)
	p1 := LogData{Data: data}
	docJSON, err := json.Marshal(p1)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %s", err)
	}

	req := esapi.IndexRequest{
		Index: index,
		//DocumentID: "2",
		Body:    strings.NewReader(string(docJSON)),
		Refresh: "true",
	}
	res, err := req.Do(context.Background(), client)
	if err != nil {
		fmt.Printf("Error indexing document: %s", err)
	}
	fmt.Println(res)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(res.Body)

	if res.IsError() {
		fmt.Printf("Failed to index document: %s", res.Status())
	} else {
		fmt.Println("Indexed document with ID: 1")
	}
	fmt.Println("ok")

}

func SendToESChan(msg *LogData) {
	ch <- msg
}

func SendTOES1() {
	for {
		select {
		case msg := <-ch:
			fmt.Println(msg.Topic, msg.Data)
			//p1 := LogData{Data: msg.Data}
			docJSON, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("Error marshalling JSON: %s", err)
			}

			req := esapi.IndexRequest{
				Index: msg.Topic,
				//DocumentID: "2",
				Body:    strings.NewReader(string(docJSON)),
				Refresh: "true",
			}
			res, err := req.Do(context.Background(), client)
			if err != nil {
				fmt.Printf("Error indexing document: %s", err)
				continue
			}
			fmt.Println(res)
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					fmt.Println(err)
				}
			}(res.Body)

			if res.IsError() {
				fmt.Printf("Failed to index document: %s", res.Status())
			} else {
				fmt.Println("Indexed document with ID: 1")
			}
			fmt.Println("ok")
		default:
			time.Sleep(time.Second)
		}

	}

}
