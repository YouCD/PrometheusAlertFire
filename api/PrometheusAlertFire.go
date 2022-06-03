package api

import (
	"PrometheusAlertFire/dist"
	"PrometheusAlertFire/pkg/config"
	"PrometheusAlertFire/pkg/log"
	"PrometheusAlertFire/pkg/metrics"
	"PrometheusAlertFire/pkg/pprof"
	"context"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
	"time"
)

//// aliyunHandler 处理阿里云 webhook 信息
//func aliyunHandler(w http.ResponseWriter, r *http.Request) {
//	if r.Method == "POST" {
//		data, err := io.ReadAll(r.Body)
//		if err != nil {
//			log.Error(err)
//		}
//		log.Debugf("alyun msg: %s", string(data))
//		aliyun := common.AliyunWebhookMsg{}
//
//		enEscapeUrl, _ := url.QueryUnescape(string(data))
//		aliyun = common.GetAliyunWebhookMsg(enEscapeUrl)
//		log.Debugf("aliyun struct: %v", aliyun)
//		fire.AlertFromAliyun(aliyun)
//		log.Info("PrometheusAlertFire receive success.")
//		return
//	}
//	fmt.Println(w, "PrometheusAlertFire is CucurbitCable alert component.")
//}

// PrometheusAlertFire 接收AlertManager webhook 信息
func PrometheusAlertFire() {

	r := mux.NewRouter()
	pprof.AttachProfiler(r)
	r.HandleFunc("/api/alert", prometheusHandler).Methods(http.MethodPost)

	receiver := r.PathPrefix("/api/receiver").Subrouter()
	// 添加 报警接收人
	receiver.HandleFunc("", createReceiver).Methods(http.MethodPost)
	// 获取 报警接收人
	receiver.HandleFunc("", listReceiver).Methods(http.MethodGet)
	// 获取 报警接收人
	receiver.HandleFunc("/search", searchReceiver).Methods(http.MethodGet)
	// 修改 报警接收人
	receiver.HandleFunc("", updateReceiver).Methods(http.MethodPut)
	// 删除 报警接收人
	receiver.HandleFunc("", delReceiver).Methods(http.MethodDelete)

	rule := r.PathPrefix("/api/subscribe").Subrouter()
	// 添加 订阅规则
	rule.HandleFunc("", createSubscribe).Methods(http.MethodPost)
	// 获取 订阅规则
	rule.HandleFunc("", listSubscribe).Methods(http.MethodGet)
	// 修改 订阅规则
	rule.HandleFunc("", updateSubscribe).Methods(http.MethodPut)
	// 删除 订阅规则
	rule.HandleFunc("", delSubscribe).Methods(http.MethodDelete)

	// prometheus
	prometheus := r.PathPrefix("/api/prometheus").Subrouter()
	prometheus.HandleFunc("/rule", fetchPrometheusRules).Methods(http.MethodGet)

	// summary
	summaryRouter := r.PathPrefix("/api/summary").Subrouter()
	summaryRouter.HandleFunc("", summary).Methods(http.MethodGet)

	//r.HandleFunc("/aliyun", aliyunHandler)
	metrics.InitMetrics()
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/logLevel", log.AtomicLevel.ServeHTTP)
	r.PathPrefix("/").Handler(http.FileServer(http.FS(dist.Dist)))

	headersOK := gorillaHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOK := gorillaHandlers.AllowedOrigins([]string{"*"})
	methodsOK := gorillaHandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"})

	handlers := gorillaHandlers.CORS(headersOK, methodsOK, originsOK)(r)
	// server
	srv := http.Server{
		Addr:    config.Cfg.Alert.ListenPort,
		Handler: handlers,
	}

	// make sure idle connections returned
	processed := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); nil != err {
			log.Panicf("Server shutdown failed, err: %v", err)
		}
		log.Info("Server gracefully shutdown")

		close(processed)
	}()
	log.Infof("Listen on http://0.0.0.0%s", config.Cfg.Alert.ListenPort)

	err := srv.ListenAndServe()
	if http.ErrServerClosed != err {
		log.Panicf("Server not gracefully shutdown, err :%v", err)
	}

	// waiting for goroutine above processed
	<-processed

}
