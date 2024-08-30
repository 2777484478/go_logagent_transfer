package conf

type AppConf struct {
	KafkaConf   `ini:"kafka"`
	TaillogConf `ini:"taillog""`
	EtcdConf    `ini:"etcd"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
	MaxSize int    `ini:"chan_max_size"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	Timeout int    `ini:"timeout"`
	Key     string `ini:"collect_log_key"`
}

type TaillogConf struct {
	Path string `ini:"path"`
}
