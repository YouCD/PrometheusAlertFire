package model

type Receiver struct {
	ID           int32  `gorm:"primaryKey;column:id;autoIncrement;not null;comment:标签值" json:"id"`
	Name         string `gorm:"column:name;type:varchar(255);not null;unique_index:name_phone;comment:名字" json:"name"`
	Telephone    string `gorm:"column:;unique;telephone;type:varchar(255);not null;unique_index:name_phone;comment:手机号"  json:"telephone"`
	Timestamp    int64  `gorm:"column:timestamp;type:bigint(11);not null;comment:修改创建时间" json:"timestamp"` // 修改时间时间
	WechatUserID string `gorm:"column:wechat_user_id;type:varchar(255);not null;comment:企业微信userid" json:"wechat_user_id"`
	Enable       int    `gorm:"column:enable;type:int(1);not null;comment:0开启1关闭" json:"enable"` // 0开启，1关闭

}

func (*Receiver) TableName() string {
	return "receiver"
}
