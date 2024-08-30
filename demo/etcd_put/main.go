package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

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

	//value := "[{\"path\":\"c:\\\\tmp\\\\nginx1.log\",\"topic\":\"web_log1\"},{\"path\":\"c:/tmp/nginx.log\",\"topic\":\"web_log\"},{\"path\":\"d:/xxx/redis.log\",\"topic\":\"redis_log\"}]"
	//value := "[{\"path\":\"./redis1.log\",\"topic\":\"redis_log\"},{\"path\":\"./nginx.log\",\"topic\":\"nginx_log\"},{\"path\":\"./mysql.log\",\"topic\":\"mysql_log\"}]"
	value := "[{\"path\":\"./nginx.log\",\"topic\":\"web_log\"},{\"path\":\"./redis.log\",\"topic\":\"redis_log\"}]"

	// 写入数据
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = cli.Put(ctx, "/logagent/192.168.18.141/collect_config", value)
	cancel()
	if err != nil {
		fmt.Println("写入数据失败:", err)
		return
	}

	fmt.Println("写入数据成功")
}
