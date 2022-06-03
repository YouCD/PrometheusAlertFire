package model

type Subscribe struct {
	ID        int32  `gorm:"primaryKey;column:id;autoIncrement;not null" json:"id"`
	RuleName  string `gorm:"column:rule_name;type:varchar(255);comment:Prometheus_rule_名称" json:"rule_name"` //Prometheus rule 名称
	Label     string `gorm:"column:label;type:varchar(255);comment:标签" json:"label"`                         // 标签
	Receiver  string `gorm:"column:receiver;type:varchar(255);comment:接收者" json:"receiver"`                  // 接受者
	Timestamp int64  `gorm:"column:timestamp;type:bigint(11);not null;comment:修改创建时间" json:"timestamp"`      // 修改创建时间
	Enable    int    `gorm:"column:enable;type:int(1);not null;comment:0开启1关闭" json:"enable"`                // 0开启，1关闭
}

func (*Subscribe) TableName() string {
	return "subscribe"
}
