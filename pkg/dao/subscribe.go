package dao

import (
	"PrometheusAlertFire/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Subscribe struct {
	table *gorm.DB
}

func NewSubscribe() *Subscribe {
	return &Subscribe{
		table: model.GetDB().Table("subscribe"),
	}
}

func (c *Subscribe) Create(obj *model.Subscribe) (err error) {
	return c.table.Create(&obj).Error
}

func (c *Subscribe) Update(obj *model.Subscribe) error {
	err := c.table.Where("id=?", obj.ID).Updates(map[string]interface{}{"rule_name": obj.RuleName, "label": obj.Label, "receiver": obj.Receiver, "enable": obj.Enable, "timestamp": obj.Timestamp}).Error
	return errors.WithMessage(err, "dao update Subscribe")
}

func (c *Subscribe) Pager(pageIndex, pageSize int) (result []*model.Subscribe, count int64) {
	res := c.table.Select("*")
	res.Count(&count)
	res.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&result)
	return
}

func (c *Subscribe) Delete(id int32) error {
	err := c.table.Delete(&model.Subscribe{}, id).Error
	return errors.WithMessage(err, "dao Delete Subscribe")
}

func (c *Subscribe) GetSubscribeByAlertname(alertname string) (result *model.Subscribe, err error) {
	err = c.table.Where("rule_name=? and enable=0", alertname).Scan(&result).Error
	return
}
func (c *Subscribe) Counter() (count int64) {
	res := c.table.Select("*")
	res.Count(&count)

	return
}
