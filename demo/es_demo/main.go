package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type Student struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Married bool   `json:"married"`
}

//func (s *Student) run() *Student {
//	fmt.Printf("%s在跑...", s.Name)
//	return s
//}

//func (s *Student) wang() {
//	fmt.Printf("%s在汪汪汪的叫...\n", s.Name)
//}

func main() {

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://172.16.1.193:9200",
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Printf("Error creating the Elasticsearch client: %s", err)
	}

	p1 := Student{Name: "Rion", Age: 22, Married: false}
	//p1.run().wang() 链式调用操作

	docJSON, err := json.Marshal(p1)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %s", err)
	}

	req := esapi.IndexRequest{
		Index:      "student",
		DocumentID: "2",
		Body:       strings.NewReader(string(docJSON)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		fmt.Printf("Error indexing document: %s", err)
	}
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
}
