package main

import (
	"fmt"
	"github.com/2777484478/log_transfer/conf"
	"github.com/2777484478/log_transfer/es"
	"github.com/2777484478/log_transfer/kafka"
	"gopkg.in/ini.v1"
)

//log  transfer
//将日志从Kafka取出，发往ES

func main() {
	//0.加载配置
	var cfg conf.LogTransferCfg
	err := ini.MapTo(&cfg, "./conf/cfg.ini")
	if err != nil {
		fmt.Printf("load ini failed :%s", err)
		return
	}
	fmt.Println(cfg)
	//1.初始化ES
	err = es.Init(cfg.EsCfg.Address)
	if err != nil {
		fmt.Println(err)
		return
	}
	//2.初始化Kafka
	//2.1 连接kafka创建消费分区
	//2.2 每个分区的消费者分别取出数据，发往ES
	err = kafka.Init([]string{cfg.KafkaCfg.Address}, cfg.KafkaCfg.Topic)
	if err != nil {
		fmt.Printf("init kafka consumer failed,err:%v\n", err)
		return
	}
	select {}

}
