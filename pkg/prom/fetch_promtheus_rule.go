package prom

import (
	"PrometheusAlertFire/pkg/config"
	"PrometheusAlertFire/pkg/log"
	"encoding/json"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"time"
)

var (
	Rules []string
)

type rule struct {
	Name string `json:"name"`
}
type group struct {
	Rules []rule `json:"rules"`
}

func fetchRules() {
	resp, err := http.Get(config.Cfg.Alert.PrometheusUrl + "/api/v1/rules?type=alert")
	if err != nil {
		log.Error(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Error(err)
	}

	result := gjson.Get(string(body), "data.groups").String()
	var groups []group
	if err = json.Unmarshal([]byte(result), &groups); err != nil {
		log.Error(err)
	}
	var tempList []string
	for _, v := range groups {
		for _, v1 := range v.Rules {
			tempList = append(tempList, v1.Name)
		}
	}
	Rules = removeDuplicateJob(tempList)
	return
}

func removeDuplicateJob(list []string) []string {
	result := make([]string, 0, len(list))
	temp := map[string]struct{}{}
	for _, item := range list {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func LoopFetchRule() {
	for {
		fetchRules()
		time.Sleep(1 * time.Hour)
	}
}
