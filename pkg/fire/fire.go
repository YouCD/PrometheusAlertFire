package fire

import (
	"PrometheusAlertFire/common"
	"PrometheusAlertFire/model"
	"PrometheusAlertFire/pkg/config"
	"PrometheusAlertFire/pkg/dao"
	"PrometheusAlertFire/pkg/log"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"hash/crc32"
	"net/url"
	"strings"

	"time"

	"github.com/pkg/errors"
)

const (
	DingTalkFiringTMPL   = "### <font color=#ff0000 size=3>%s </font>  \n\n**开始时间:** %s \n\n**结束时间:** %s \n\n**故障主机IP:** %s \n\n**告警级别:** %s \n\n %s \n\n[Prometheus](%s)"
	ReDingTalkFiringTMPL = "### <font color=#00C957 size=3>%s </font>  \n\n**开始时间:** %s \n\n**结束时间:** %s \n\n**故障主机IP:** %s \n\n**告警级别:** %s \n\n %s \n\n[Prometheus](%s)"

	wxTextTMPL = "[%s]()\n>**[%s]()**\n>`告警级别:`%s\n`开始时间:`%s\n`结束时间:`%s\n`故障主机IP:` %s\n**%s**\n[Prometheus](%s)\n"
)

// AlertFromPrometheus 发送告警
func AlertFromPrometheus(prometheusMsg common.Prometheus) {
	for _, Msg := range prometheusMsg.Alerts {
		//格式化开始时间
		at, err := timeFormat(Msg.StartsAt)
		if err != nil {
			log.Warn(err)
		}
		//格式化结束时间
		et, err := timeFormat(Msg.EndsAt)
		if err != nil {
			log.Warn(err)
		}
		// 处理 prometheus的 url
		prometheusUrl := formatPrometheusUrl(Msg.GeneratorUrl)

		switch {
		case Msg.Status == "resolved":

			wxText := fmt.Sprintf(wxTextTMPL, config.Cfg.Alert.Title, Msg.Labels["alertname"]+"    [Resolved]", Msg.Labels["severity"], at, et, Msg.Labels["instance"], Msg.Annotations.Description, prometheusUrl)
			DingTalkText := fmt.Sprintf(ReDingTalkFiringTMPL, config.Cfg.Alert.Title+Msg.Labels["alertname"]+"    [Resolved]", at, et, Msg.Labels["instance"], Msg.Labels["severity"], Msg.Annotations.Description, prometheusUrl)
			title := Msg.Labels["alertname"]

			adapter := NewReceiverAdapter(title, DingTalkText, wxText, Msg.Labels["alertname"], nil, nil)
			if adapter == nil {
				return
			}
			err := adapter.FireMsg()
			if err != nil {
				log.Warn(err)
			}
			log.Infof("%s is resolved.", Msg.Labels["alertname"])

		case Msg.Status == "firing":

			// 处理各个通道的消息模板
			wxText := fmt.Sprintf(wxTextTMPL, config.Cfg.Alert.Title, Msg.Labels["alertname"], Msg.Labels["severity"], at, et, Msg.Labels["instance"], Msg.Annotations.Description, prometheusUrl)
			DingTalkText := fmt.Sprintf(DingTalkFiringTMPL, Msg.Labels["alertname"], at, et, Msg.Labels["instance"], Msg.Labels["severity"], Msg.Annotations.Description, prometheusUrl)
			title := config.Cfg.Alert.Title + Msg.Labels["alertname"]

			MergeAlert(title, DingTalkText, wxText, &Msg)

		}

	}

}

// formatPrometheusUrl 格式化 Prometheus 接口地址
func formatPrometheusUrl(u string) string {
	OldParse, err := url.Parse(u)
	if err != nil {
		return u
	}
	NewParse, err := url.Parse(config.Cfg.Alert.PrometheusUrl)
	if err != nil {
		log.Warnf("Config Alert.PrometheusUrl field is err: %s", err.Error())
		return u
	}

	OldParse.Host = NewParse.Host
	return OldParse.String()
}

// MergeAlert 合并告警信息
func MergeAlert(title, dingTalkText, wxText string, alert *common.Alert) {
	// 1. 判断此条alert是否在规则列表中
	_, receivers, flag := checkInRule(alert)

	// 2. 获取 手机号 企业微信userid 列表
	var Mobiles []string
	var AtMobiles string
	var wxUserIDs []string
	for _, v := range receivers {
		AtMobiles += "@" + v.Telephone + " "
		Mobiles = append(Mobiles, v.Telephone)
		wxUserIDs = append(wxUserIDs, v.WechatUserID)
	}

	// 3. 发送消息
	switch {
	// 3.1. 在规则列表中 @指定人
	case flag:
		DingTalkAtText := dingTalkText + "\n\n" + AtMobiles

		adapter := NewReceiverAdapter(title, DingTalkAtText, wxText, alert.Labels["alertname"], Mobiles, wxUserIDs)
		if adapter == nil {
			return
		}
		err := adapter.FireMsg()
		if err != nil {
			log.Warn(err)
			return
		} else {
			goto LOG
		}
	// 3.2. 不在规则列表中 就走默认通道 不 @具体的某个人
	default:
		adapter := NewReceiverAdapter(title, dingTalkText, wxText, alert.Labels["alertname"], nil, nil)
		if adapter == nil {
			return
		}
		err := adapter.FireMsg()
		if err != nil {
			log.Warn(err)
			return
		} else {
			goto LOG
		}
	}

LOG:
	log.Infof("%s is firing.", alert.Labels["alertname"])
}

func parseLabel(label string) (labels map[string]string) {
	StrList := strings.Split(label, ",")
	labels = make(map[string]string)
	for _, v := range StrList {
		lv := strings.Split(v, "=")
		labels[lv[0]] = lv[1]
	}
	return
}

func isMapSubset[K, V comparable](m, sub map[K]V) bool {
	if len(sub) > len(m) {
		return false
	}
	for k, vsub := range sub {
		if vm, found := m[k]; !found || vm != vsub {
			return false
		}
	}
	return true
}

// checkInRule 检查alert是否在规则列表中
func checkInRule(alert *common.Alert) (hs int, owners []*model.Receiver, flag bool) {
	result, err := dao.NewSubscribe().GetSubscribeByAlertname(alert.Labels["alertname"])
	if err != nil {
		log.Error(err)
		return
	}
	if result == nil {
		return 0, nil, false
	}

	var Receivers []*model.Receiver
	array := gjson.Parse(result.Receiver).Array()
	for _, v := range array {
		i := gjson.Get(v.String(), "id").Int()
		Receiver, err := dao.NewReceiver().GetReceiverByID(int(i))
		if err != nil {
			log.Error(err)
			return
		}
		if Receiver != nil {
			Receivers = append(Receivers, Receiver)
		}

	}
	hs = checksum(alert)
	labelMap := parseLabel(result.Label)
	if result.RuleName == alert.Labels["alertname"] {
		flag = isMapSubset(alert.Labels, labelMap)
		return hs, Receivers, flag
	}
	return 0, nil, false
}

// checksum 将alert信息hash为int
func checksum(data interface{}) int {
	switch data.(type) {
	case common.Alert:
		jsb, _ := json.Marshal(data)
		v := int(crc32.ChecksumIEEE(jsb))
		if v >= 0 {
			return v
		}
		if -v >= 0 {
			return -v
		}
	case string:
		v := int(crc32.ChecksumIEEE([]byte(data.(string))))
		if v >= 0 {
			return v
		}
		if -v >= 0 {
			return -v
		}
	}
	return 0
}

// timeFormat 格式化时间
func timeFormat(timeStr string) (transTime string, err error) {
	if timeStr == "0001-01-01T00:00:00Z" {
		return "请即时处理故障", nil
	}
	//2006-01-02 15:04:05 是golang的时间模板，据说是golang语言的诞生时间，2006-01-02 15:04:05类似于我们熟悉的YYYY-MM-dd HH:mm:ss
	result, err := time.ParseInLocation(time.RFC3339Nano, timeStr, time.Local)
	if err != nil {
		return timeStr, errors.WithMessage(err, "timeFormat")
	}
	timeLayout := "2006-01-02 15:04:05"
	transTime = result.Local().Format(timeLayout)
	return
}

//func AlertFromAliyun(msg common.AliyunWebhookMsg) {
//
//	aliyunRul := "https://signin.aliyun.com/daddylab.onaliyun.com/login.htm#/main"
//	titleend := "故障告警信息"
//	wxtext := "[" + "阿里云监控:" + titleend + "](" + aliyunRul + ")\n>**" + msg.AlertName + "**\n>`告警级别:`" + "info" + "\n`持续时间:`" + msg.LastTime + "\n`故障实例:`" + msg.InstanceName + "\n"
//
//	metrics.FireWXAlertCounter.Inc()
//	err := MsgHandler(wx, wxtext, "")
//	if err != nil && errors.Is(err, ErrReqLimit) {
//		log.Warnf("WX: FireAlert failed. reason: ReqLimit. %s: %s  ", "InstanceName", msg.InstanceName)
//		return
//	} else if err != nil {
//		log.Errorf("WX: FireAlert failed. %s: %s  err:%s", "app", "InstanceName", msg.InstanceName, err.Error())
//		return
//	}
//	metrics.FireWXAlertSuccess.Inc()
//}

//func ClearNotifyList() {
//	if len(common.IsNotifyList) > 0 {
//		common.IsNotifyList = nil
//		log.Info("AlertFiler:  NotifyList is cleared")
//	}
//}
func isNowInTimeRange(startTimeStr, endTimeStr string) bool {
	//当前时间
	now := time.Now()
	//当前时间转换为"年-月-日"的格式
	format := now.Format("2006-01-02")
	//转换为time类型需要的格式
	layout := "2006-01-02 15:04"
	//将开始时间拼接“年-月-日 ”转换为time类型
	timeStart, _ := time.ParseInLocation(layout, format+" "+startTimeStr, time.Local)

	//将结束时间拼接“年-月-日 ”转换为time类型
	timeEnd, _ := time.ParseInLocation(layout, format+" "+endTimeStr, time.Local)

	//使用time的Before和After方法，判断当前时间是否在参数的时间范围
	return now.Before(timeEnd) && now.After(timeStart)
}
