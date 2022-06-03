package fire

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dyvmsapi"
	"github.com/pkg/errors"
	"strings"
)

var (
	CallLimitERR = errors.New("该被叫号触发被叫流控")
)

type AliYunPhone struct {
	AccessKeyId  string
	AccessSecret string
	TTSCode      string
}

//func GetAliYunPhone() *AliYunPhone {
//	return &AliYunPhone{
//		AccessKeyId:  config.Cfg.AliYun.AccessKeyId,
//		AccessSecret: config.Cfg.AliYun.AccessSecret,
//		TTSCode:      config.Cfg.AliYun.TTSCode,
//	}
//
//}

func (m *AliYunPhone) FireMsg(Messages string, PhoneNumbers string) (err error) {
	mobiles := strings.Split(PhoneNumbers, ",")
	for _, mobile := range mobiles {
		client, err := dyvmsapi.NewClientWithAccessKey("cn-hangzhou", m.AccessKeyId, m.AccessSecret)
		if err != nil {
			return err
		}
		request := dyvmsapi.CreateSingleCallByTtsRequest()
		request.Scheme = "https"
		//阿里云电话被叫显号，必须是已购买的号码
		//request.CalledShowNumber = CalledShowNumber
		request.CalledNumber = mobile
		request.TtsCode = m.TTSCode
		request.TtsParam = `{"code":"` + Messages + `"}`
		request.PlayTimes = requests.NewInteger(2)

		response, err := client.SingleCallByTts(request)
		if err != nil {
			return err
		}
		if strings.Contains(response.Message, "该被叫号触发被叫流控") {
			return CallLimitERR
		}
	}
	return nil
}
