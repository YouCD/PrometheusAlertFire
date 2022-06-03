package api

import (
	"PrometheusAlertFire/common"
	"PrometheusAlertFire/pkg/fire"
	"PrometheusAlertFire/pkg/log"
	"PrometheusAlertFire/pkg/prom"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// prometheusHandler 处理AlertManager webhook 信息
func prometheusHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
	}
	alert := common.Prometheus{}
	err = json.Unmarshal(data, &alert)
	if err != nil {
		log.Error(err)
	}
	log.Debug(string(data))
	fire.AlertFromPrometheus(alert)
	log.Info("Receive success.")
	return
}
func fetchPrometheusRules(w http.ResponseWriter, r *http.Request) {
	var a = struct {
		Data interface{} `json:"data"`
	}{
		Data: prom.Rules,
	}
	fmt.Fprint(w, successResponse("获取成功", true, a))
	return
}
