package main

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var Ksbdl = "BAODELU"

func main() {
	// 建立连接
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"}, // 替换为您的etcd实例地址
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("连接etcd失败:", err)
		return
	}
	defer cli.Close()

	// 写入数据
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = cli.Put(ctx, Ksbdl, "DSB")
	cancel()
	if err != nil {
		fmt.Println("写入数据失败:", err)
		return
	}

	fmt.Println("写入数据成功", Ksbdl)

	// 读取数据
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := cli.Get(ctx, "foo")
	cancel()
	if err != nil {
		fmt.Println("读取数据失败:", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("获取数据 %s : %s\n", ev.Key, ev.Value)
	}
}
