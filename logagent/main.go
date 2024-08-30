package main

import (
	"fmt"
	"github.com/2777484478/logagent/conf"
	"github.com/2777484478/logagent/etcd"
	"github.com/2777484478/logagent/kafka"
	"github.com/2777484478/logagent/taillog"
	"github.com/2777484478/logagent/utils"
	"gopkg.in/ini.v1"
	"sync"
	"time"
)

var cfg = new(conf.AppConf)

//var cfg = &conf.AppConf{}

func main() {

	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Printf("load ini failed err :%v\n", err)
		return
	}
	//1.初始化kafka连接
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.MaxSize)

	if err != nil {
		fmt.Println("kafka初始化连接失败：", err)
		return
	}
	fmt.Println("Init Kafka Success!")

	//2.初始化etcd

	err = etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Printf("connect etcd failed : %s,\n", err)
		return
	}
	fmt.Println("connect etcd success!!!")
	//2.1从etcd中获取收集项的配置
	ipStr, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	etcdConfKey := fmt.Sprintf(cfg.EtcdConf.Key, ipStr)
	fmt.Println("ldc", etcdConfKey)

	logEntryConf, err1 := etcd.GetConf(etcdConfKey)
	if err != nil {
		fmt.Printf(" etcd getConf failed : %s,\n", err1)
		return
	}
	fmt.Printf("Get etcd conf  sucess !!! value: %v \n", logEntryConf)

	//2.2派一个哨兵监视日志项的变化，有变化及时通知我的logAgent实现热加载

	for index, value := range logEntryConf {
		fmt.Println(index, value)
	}
	//3. 收集日志发往 Kafka
	taillog.Init(logEntryConf)

	for index, value := range logEntryConf {
		fmt.Printf("index:%v,value:%v\n", index, value)
	}
	newConfChan := taillog.NewConfChan()
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(etcdConfKey, newConfChan)
	wg.Wait()
	//select {}
}
