package taillog

import (
	"fmt"
	"github.com/2777484478/logagent/etcd"
	"time"
)

var tskMgr *tailLogMgr

// tailTask 管理者
type tailLogMgr struct {
	logEntry    []*etcd.LogEntry
	tskMap      map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}

func Init(logEntryConf []*etcd.LogEntry) {
	tskMgr = &tailLogMgr{
		logEntry:    logEntryConf, //把当前日志收集项 配置信息保存起来
		tskMap:      make(map[string]*TailTask, 16),
		newConfChan: make(chan []*etcd.LogEntry), //无缓冲区通道
	}
	for _, logEntry := range logEntryConf {
		//初始化时候起了多少个 TailTask 都要记录下来，方便后续判断
		mk := fmt.Sprintf("%s_%s", logEntry.Path, logEntry.Topic)
		tailObj := NewTailTask(logEntry.Path, logEntry.Topic)
		tskMgr.tskMap[mk] = tailObj
	}
	go tskMgr.run()
}

func (t *tailLogMgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			//1.配置新增
			//2.删除
			//3.更改操作
			for _, conf := range newConf {
				mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				_, ok := t.tskMap[mk]
				if ok {
					continue
				} else {
					tailObj := NewTailTask(conf.Path, conf.Topic)
					t.tskMap[mk] = tailObj
				}
				for _, c1 := range t.logEntry {
					isDelete := true
					for _, c2 := range newConf {
						if c2.Path == c1.Path && c2.Topic == c1.Topic {
							isDelete = false
							continue
						}
					}
					if isDelete {
						mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
						t.tskMap[mk].cancelFunc()
					}
				}
			}
			fmt.Println("新配置变更", newConf)
		default:
			time.Sleep(time.Second)
		}
	}
}

func NewConfChan() chan<- []*etcd.LogEntry {
	return tskMgr.newConfChan
}
