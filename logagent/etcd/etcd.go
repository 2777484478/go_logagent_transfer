package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var client *clientv3.Client

type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

//var a *logEntry  创建指针类型变量
//var b = &logEntry{} 创建对象指针
//var c = new(logEntry) 创建对象指针

func Init(addr string, timeout time.Duration) (err error) {

	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: timeout,
	})
	if err != nil {
		fmt.Printf("connect etcd failed : %s,\n", err)
		return
	}
	//defer client.Close()
	//这里不能使用 defer 关闭连接 否则会报错
	//rpc error: code = Canceled desc = grpc: the client connection is closing
	return
}

func GetConf(key string) (logEntryConf []*LogEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := client.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Println("读取数据失败:", err)
		return
	}
	fmt.Println(resp)
	for _, ev := range resp.Kvs {
		fmt.Printf("获取数据key: %s : value:%s\n", ev.Key, ev.Value)
		err = json.Unmarshal(ev.Value, &logEntryConf)
		if err != nil {
			fmt.Printf("Unmarshal etcd failed \n,err:%v\n", err)
			return
		}
	}
	return
}

func WatchConf(key string, newConfCh chan<- []*LogEntry) {
	ch := client.Watch(context.Background(), key)
	fmt.Println("Watching...")

	for wresp := range ch {
		for _, evt := range wresp.Events {
			fmt.Printf("Tyep:%v, key:%v, value: %v \n ", evt.Type, string(evt.Kv.Key), string(evt.Kv.Value))
			//通知别人
			var newConf []*LogEntry
			if evt.Type != clientv3.EventTypeDelete {
				//如果不是删除操作
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					fmt.Printf("json failed ,err : %v\n", err)
					continue
				}
			}
			fmt.Printf("get new conf:%v\n", newConf)
			newConfCh <- newConf
		}
	}
}
