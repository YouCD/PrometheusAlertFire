package common

type Annotations struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  Annotations       `json:"annotations"`
	StartsAt     string            `json:"startsAt"`
	EndsAt       string            `json:"endsAt"`
	GeneratorUrl string            `json:"generatorURL"` //prometheus 告警返回地址
}

type Prometheus struct {
	Status      string
	Alerts      []Alert
	Externalurl string `json:"externalURL"` //alertmanage 返回地址
}
