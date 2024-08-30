package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {

	client, err := clientv3.New(clientv3.Config{Endpoints: []string{"http://localhost:2379"}, DialTimeout: 5 * time.Second})
	if err != nil {
		fmt.Printf("connect etcd failed : %s,\n", err)
		return
	}
	defer client.Close()
	fmt.Println("connect etcd success!!!")

	//watch

	ch := client.Watch(context.Background(), "ldc")
	fmt.Println("Watching...")

	for wresp := range ch {
		for _, evt := range wresp.Events {
			fmt.Printf("Tyep:%v, key:%v, value: %v \n ", evt.Type, string(evt.Kv.Key), string(evt.Kv.Value))
		}
	}

}
