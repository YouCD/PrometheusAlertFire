package fire

import (
	"PrometheusAlertFire/pkg/config"
	"PrometheusAlertFire/pkg/log"
	"PrometheusAlertFire/pkg/metrics"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	ErrReqLimit = errors.New("api freq out of limit")
)

const WorkWechatWebhookPrefix = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="

type mark struct {
	Content             string   `json:"content"`
	MentionedMobileList []string `json:"mentioned_mobile_list"`
}
type WXMessage struct {
	Msgtype  string `json:"msgtype"`
	Markdown mark   `json:"markdown"`
}

type WechatBoot struct {
	WXMessage     WXMessage
	WechatBootUrl string
}

func NewWechatBoot(content string, atUsers []string) *WechatBoot {
	if len(atUsers) > 0 {
		idtext := ""
		for _, id := range atUsers {
			idtext += "<@" + id + ">"
		}
		content += idtext
	}

	return &WechatBoot{
		WXMessage: WXMessage{
			Msgtype: "markdown",
			Markdown: mark{
				Content:             content,
				MentionedMobileList: atUsers,
			},
		},
		WechatBootUrl: WorkWechatWebhookPrefix + config.Cfg.WorkWechat.Key,
	}

}

func (m *WechatBoot) FireMsg() (err error) {
	metrics.FireWXAlertCounter.Inc()

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(m.WXMessage)
	if err != nil {
		return err
	}
	var tr *http.Transport
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	res, err := client.Post(m.WechatBootUrl, "application/json", b)
	if err != nil {
		log.Error(err)
		return errors.WithMessage(err, "PostToWeiXin")
	}
	if res == nil {
		return errors.WithMessage(err, "PostToWeiXin response is nil")
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if strings.Contains(string(result), "api freq out of limit") {
		return ErrReqLimit
	}
	if gjson.Get(string(result), "errcode").Int() != 0 {
		return errors.New(string(result))
	}
	metrics.FireWXAlertSuccess.Inc()

	return nil
}
