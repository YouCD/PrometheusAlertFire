package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

//type AliYun struct {
//	AccessSecret string
//	AccessKeyId  string
//	TTSCode      string
//}

const (
	WorkWechatNotify = iota + 1
	DingTalkNotify
)

type dingTalk struct {
	Sign  string
	Token string
}

type alert struct {
	Title         string
	PrometheusUrl string
	RuleName      string
	ListenPort    string
	LogLevel      string
	Enabled       bool
	DefaultNotify string
}
type workWechat struct {
	Key string
}
type mysql struct {
	HostAndPort string
	DBName      string
	Password    string
	User        string
}

type silences struct {
	StartTimeStr string
	EndTimeStr   string
	Enabled      bool
}

type config struct {
	Alert      alert
	Mysql      mysql
	DingTalk   dingTalk
	WorkWechat workWechat
	Silences   silences
}

var (
	Cfg           = new(config)
	ConfigureFile string
)

func init() {
	err := newConfig()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func newConfig() (err error) {
	v := viper.New()

	//WxUrl := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
	v.AddConfigPath("../conf") //设置读取的文件路径
	v.AddConfigPath("../")     //设置读取的文件路径
	v.AddConfigPath("./")      //设置读取的文件路径
	v.SetConfigName("config")  //设置读取的文件名
	v.SetConfigType("yaml")    //设置文件的类型
	err = v.ReadInConfig()
	if err != nil {
		return err
	}
	ConfigureFile = v.ConfigFileUsed()

	var dt dingTalk
	err = v.UnmarshalKey("DingTalk", &dt)
	if err != nil {
		return err
	}
	Cfg.DingTalk = dt

	var at alert
	err = v.UnmarshalKey("Alert", &at)
	if err != nil {
		return err
	}
	Cfg.Alert = at
	p := v.GetString("Alert.ListenPort")
	if p == "" {
		p = "8080"
	}
	Cfg.Alert.ListenPort = ":" + p

	var m mysql
	err = v.UnmarshalKey("MySQL", &m)
	if err != nil {
		return err
	}
	Cfg.Mysql = m

	var w workWechat
	err = v.UnmarshalKey("WorkWechat", &w)
	if err != nil {
		return err
	}
	Cfg.WorkWechat = w

	if v.IsSet("Silences.TimeArray") {
		TimeArray := v.GetString("Silences.TimeArray")
		if !strings.Contains(TimeArray, "~") {
			fmt.Println("TimeArray 格式不正确")
			os.Exit(-1)
		}

		s := strings.Split(TimeArray, "~")
		Cfg.Silences.StartTimeStr = s[0]
		Cfg.Silences.EndTimeStr = s[1]
	}
	Cfg.Silences.Enabled = v.GetBool("Silences.Enabled")

	return err
}

func NotifyType() int {
	switch {
	case strings.ToLower(Cfg.Alert.DefaultNotify) == "workwechat":
		return WorkWechatNotify
	case strings.ToLower(Cfg.Alert.DefaultNotify) == "dingtalk":
		return DingTalkNotify
	default:
		fmt.Println("请配置默认的通知方式 DefaultNotify")
		os.Exit(1)
	}
	return 0
}
