package main

import (
	"PrometheusAlertFire/api"
	"PrometheusAlertFire/pkg/config"
	"PrometheusAlertFire/pkg/log"
	"PrometheusAlertFire/pkg/prom"
	"time"
)

func main() {
	log.Debugf("Config file is %s", config.ConfigureFile)
	go func() {
		for {
			time.Sleep(time.Hour * 24)
		}
	}()
	go prom.LoopFetchRule()
	api.PrometheusAlertFire()
}
