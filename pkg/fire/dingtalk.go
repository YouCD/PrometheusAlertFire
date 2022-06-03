package fire

import (
	"PrometheusAlertFire/pkg/config"
	"PrometheusAlertFire/pkg/metrics"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"strings"

	"net/url"
	"time"

	"net/http"
)

const DingTalkWebhookPrefix = "https://oapi.dingtalk.com/robot/send?access_token="

type Mark struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	//IsAtAll       bool     `json:"isAtAll"`
	//AtDingtalkIds []string `json:"atDingtalkIds"`
}
type Message struct {
	Msgtype  string `json:"msgtype"`
	Markdown Mark   `json:"markdown"`
	At       At     `json:"at"`
}
type DingTalkBoot struct {
	Message
	bootUrl string
}

func NewDingTalkBoot(title, text string, mobiles []string) *DingTalkBoot {
	var at At
	if mobiles == nil {
		at.AtMobiles = []string{}
	} else {
		at.AtMobiles = mobiles
	}
	return &DingTalkBoot{
		Message: Message{
			Msgtype: "markdown",
			Markdown: Mark{
				Title: title,
				Text:  text,
			},
			At: at,
		},
		bootUrl: DingTalkWebhookPrefix + config.Cfg.DingTalk.Token,
	}
}

func sign() (timestamp int64, sign string) {
	timestamp = time.Now().UnixMilli()
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, config.Cfg.DingTalk.Sign)
	hash := hmac.New(sha256.New, []byte(config.Cfg.DingTalk.Sign))
	hash.Write([]byte(stringToSign))
	signData := hash.Sum(nil)
	return timestamp, url.QueryEscape(base64.StdEncoding.EncodeToString(signData))
}

func (d *DingTalkBoot) FireMsg() (err error) {
	metrics.FireDingTalkAlertCounter.Inc()
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(d.Message)
	timestamp, sign := sign()
	res, err := http.Post(fmt.Sprintf("%s&timestamp=%d&sign=%s", d.bootUrl, timestamp, sign), "application/json", b)
	if err != nil {
		return errors.WithMessage(err, "PostToDingTalk")
	}
	if res == nil {
		return errors.WithMessage(err, "PostToDingTalk response is nil")
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if strings.Contains(string(result), "send too fast") {
		return ErrReqLimit
	}
	if gjson.Get(string(result), "errcode").Int() != 0 {
		return errors.New(string(result))
	}
	metrics.FireDingTalkAlertSuccess.Inc()
	return nil

}
