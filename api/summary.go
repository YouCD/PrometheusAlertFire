package api

import (
	"PrometheusAlertFire/pkg/config"
	"PrometheusAlertFire/pkg/dao"
	"PrometheusAlertFire/pkg/metrics"
	"fmt"
	"net/http"
)

type summaryInfo struct {
	FireTotal        float64 `json:"fireTotal"`
	FireSuccessTotal float64 `json:"fireSuccessTotal"`
	SubscribeTotal   int     `json:"subscribeTotal"`
	ReceiverTotal    int     `json:"receiverTotal"`
}

func summary(w http.ResponseWriter, r *http.Request) {
	var sm summaryInfo
	sm.ReceiverTotal = int(dao.NewReceiver().Counter())
	sm.SubscribeTotal = int(dao.NewSubscribe().Counter())
	switch {
	// 钉钉通道
	case config.NotifyType() == config.DingTalkNotify:
		sm.FireSuccessTotal = metrics.GetCounterValue(metrics.FireDingTalkAlertSuccess)
		sm.FireTotal = metrics.GetCounterValue(metrics.FireDingTalkAlertCounter)
	case config.NotifyType() == config.WorkWechatNotify:
		sm.FireSuccessTotal = metrics.GetCounterValue(metrics.FireWXAlertSuccess)
		sm.FireTotal = metrics.GetCounterValue(metrics.FireWXAlertCounter)
	}

	//TmpPageIndex := r.FormValue("page_index")
	//PageIndex, err := strconv.Atoi(TmpPageIndex)
	//if err != nil {
	//	log.Warn(err)
	//	fmt.Fprint(w, successResponse("请指定pageIndex", false, err))
	//	return
	//}
	//
	//TmpPageSize := r.FormValue("page_size")
	//PageSize, err := strconv.Atoi(TmpPageSize)
	//if err != nil {
	//	log.Warn(err)
	//	fmt.Fprint(w, successResponse("请指定pageIndex", false, err))
	//	return
	//}
	//
	//result, count := dao.NewReceiver().Pager(PageIndex, PageSize)
	//
	var a = struct {
		Data interface{} `json:"data"`
	}{
		sm,
	}

	fmt.Fprint(w, successResponse("获取成功", true, a))
	return
}
