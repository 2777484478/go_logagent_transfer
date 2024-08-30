package taillog

import (
	"context"
	"fmt"
	"github.com/2777484478/logagent/kafka"
	"github.com/hpcloud/tail"
)

var tailObj *tail.Tail

// TailTask 一个日志收集任务
type TailTask struct {
	path       string
	topic      string
	instance   *tail.Tail
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewTailTask(path, topic string) (tailObj *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		path:       path,
		topic:      topic,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	tailObj.init()
	return
}

func (t *TailTask) init() {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true}
	var err error
	t.instance, err = tail.TailFile(t.path, config)
	fmt.Println(t.path)
	if err != nil {
		fmt.Println("tail file failed", err)
	}
	go t.run() //直接采集数据发送到Kafka
}

func (t *TailTask) run() {
	for {
		select {
		case <-t.ctx.Done():
			fmt.Printf("tail task:%s_%s 结束了 \n", t.path, t.topic)

		case line := <-t.instance.Lines: // 从 tailObj 通道中一行行的读取日志数据
			//3.2 发往Kafka
			//先把日志发往一个通道中
			fmt.Printf("get log data from %s success, log:%v\n", t.path, line.Text)
			kafka.SendToChan(t.topic, line.Text)
		}
	}
}
