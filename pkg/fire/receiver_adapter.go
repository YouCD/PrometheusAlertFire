package fire

import (
	"PrometheusAlertFire/pkg/config"
	"PrometheusAlertFire/pkg/log"
)

type ReceiverAdapter interface {
	FireMsg() error
}

func NewReceiverAdapter(title, DingTalkAtText, wxText, alertName string, mobiles, wxUserIDs []string) (Receiver ReceiverAdapter) {

	switch {
	case config.Cfg.Silences.Enabled:
		inTimeArray := isNowInTimeRange(config.Cfg.Silences.StartTimeStr, config.Cfg.Silences.EndTimeStr)
		if inTimeArray && len(mobiles) == 0 {
			log.Warnf("Silences timeArray: [ %s ~ %s ], %s is block", config.Cfg.Silences.StartTimeStr, config.Cfg.Silences.EndTimeStr, alertName)
			return
		} else if config.Cfg.Alert.Enabled {
			switch {
			// 钉钉通道
			case config.NotifyType() == config.DingTalkNotify:
				return NewDingTalkBoot(title, DingTalkAtText, mobiles)
			case config.NotifyType() == config.WorkWechatNotify:
				return NewWechatBoot(wxText, wxUserIDs)
			}
		}
	case config.Cfg.Alert.Enabled:
		switch {
		// 钉钉通道
		case config.NotifyType() == config.DingTalkNotify:
			return NewDingTalkBoot(title, DingTalkAtText, mobiles)
		case config.NotifyType() == config.WorkWechatNotify:
			return NewWechatBoot(wxText, wxUserIDs)
		}
	}

	return nil
}
