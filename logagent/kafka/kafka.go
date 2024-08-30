package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"time"
)

type logData struct {
	topic string
	data  string
}

var (
	client      sarama.SyncProducer
	logDataChan chan *logData
)

func Init(addr []string, maxSize int) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	client, err = sarama.NewSyncProducer(addr, config)

	if err != nil {
		fmt.Println("producer closed , err:", err)
		return
	}
	//初始化logDataChan
	logDataChan = make(chan *logData, maxSize)
	//开启后台 goroutine 从通道取数据
	go sendToKafka()
	return
}

func SendToChan(topic string, data string) {
	msg := &logData{topic: topic, data: data}
	logDataChan <- msg
}

// 往 kafka 发送日志函数
func sendToKafka() {
	for {
		select {
		case ld := <-logDataChan:
			msg := &sarama.ProducerMessage{
				Topic: ld.topic,
				Value: sarama.StringEncoder(ld.data),
			}
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg failed , err :", err)
				return
			}
			fmt.Printf("pid : %v offset : %v\n\n", pid, offset)
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}
