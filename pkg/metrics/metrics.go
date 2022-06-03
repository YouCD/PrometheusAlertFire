package metrics

import (
	"PrometheusAlertFire/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

var (
	FireWXAlertCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "fire_alert_to_wx_total",
			Help: "The total number of fire alert to wechat events",
		},
	)
	FireWXAlertSuccess = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "fire_alert_to_wx_success",
			Help: "The total number of fire alert to wechat success events",
		},
	)

	//FireAliYunAlertCounter = prometheus.NewCounter(
	//	prometheus.CounterOpts{
	//		Name: "fire_alert_to_aliyun_total",
	//		Help: "The total number of fire alert to aliyun events",
	//	},
	//)
	//
	//FireAliYunAlertSuccess = prometheus.NewCounter(
	//	prometheus.CounterOpts{
	//		Name: "fire_alert_to_aliyun_success",
	//		Help: "The total number of fire alert to aliyun success events",
	//	},
	//)

	FireDingTalkAlertCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "fire_alert_to_dingtalk_total",
			Help: "The total number of fire alert to dingtalk events",
		},
	)

	FireDingTalkAlertSuccess = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "fire_alert_to_dingtalk_success",
			Help: "The total number of fire alert to dingtalk success events",
		},
	)
)

func InitMetrics() {
	switch {
	// 钉钉通道
	case config.NotifyType() == config.DingTalkNotify:
		prometheus.MustRegister(FireDingTalkAlertSuccess, FireDingTalkAlertCounter)
	case config.NotifyType() == config.WorkWechatNotify:
		prometheus.MustRegister(FireWXAlertCounter, FireWXAlertSuccess)

	}

}

func GetCounterValue(m prometheus.Counter) float64 {
	pb := &dto.Metric{}
	m.Write(pb)
	return pb.GetCounter().GetValue()
}
