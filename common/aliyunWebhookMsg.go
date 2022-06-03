package common

import (
	"strings"
)

type MetricColumns struct {
	InstanceId string  `json:"instanceId"`
	Average    float32 `json:"average"`
	Vip        string  `json:"vip"`
	Timestamp  int     `json:"timestamp"`
}

type Message struct {
	MetricColumns MetricColumns `json:"metricColumns"`
	RuleName      string        `json:"ruleName"`
	Time          int           `json:"time"`
}


type AliyunWebhookMsg struct {
	AlertName    string `json:"alertName"`    //规则名称
	CurValue     string `json:"curValue"`     //当前值
	InstanceName string `json:"instanceName"` //实例名称
	LastTime     string `json:"lastTime"`     //持续时间
}

func GetAliyunWebhookMsg(aliMSg string) (msg AliyunWebhookMsg) {
	l := strings.Split(aliMSg, "&")
	for _, v := range l {
		kv := strings.Split(v, "=")
		switch {
		case kv[0] == "alertName":
			msg.AlertName = kv[1]
		case kv[0] == "curValue":
			msg.CurValue = kv[1]
		case kv[0] == "instanceName":
			msg.InstanceName = kv[1]
		case kv[0] == "lastTime":
			msg.LastTime = kv[1]
		}
	}

	return
}
