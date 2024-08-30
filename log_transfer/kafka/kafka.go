package kafka

import (
	"fmt"
	"github.com/2777484478/log_transfer/es"
	"github.com/IBM/sarama"
)

func Init(addrs []string, topic string) (err error) {
	consumer, err := sarama.NewConsumer(addrs, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Printf("连接kafaka：分区列表：%v", partitionList)
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return err
		}
		//defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				//直接发给ES
				fmt.Println("发送到es")
				//es.SendTOES(topic, string(msg.Value))函数调用函数
				//	优化一下，直接放到通道中
				ld := es.LogData{Topic: topic,
					Data: string(msg.Value)}
				es.SendToESChan(&ld)
			}
		}(pc)
	}
	return err
}
