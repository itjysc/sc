package common

type Metric struct {
	Metric    string   `json:"metric"`    //“Metric”定义指标的名称，如cpu,mem等等。后面的“ `json:"metric"` ”表示序列化json中的key.
	Endpoint  string   `json:"endpoint"`  //“Endpoint”定义主机名
	Tag       []string `json:"tag"`       //“Tag”打标签，可以用来识别当前的操作系统
	Value     float64  `json:"value"`     //“Value”监控指标当前的值
	Timestamp int64    `json:"timestamp"` //“Timestamp”定义时间戳，用来标记你的值是合适传来的
}
