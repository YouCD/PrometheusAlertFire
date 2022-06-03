package dao

import (
	"PrometheusAlertFire/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

var (
	ErrIsNull = errors.New("Value is null")
)

type Receiver struct {
	table *gorm.DB
}

func NewReceiver() *Receiver {
	return &Receiver{
		table: model.GetDB().Table("receiver"),
	}
}

func (c *Receiver) Create(obj *model.Receiver) (err error) {
	if obj.Name == "" || obj.Telephone == "" {
		return ErrIsNull
	}
	obj.Timestamp = time.Now().Unix()
	return c.table.Create(&obj).Error
}

func (c *Receiver) Update(obj *model.Receiver) error {
	err := c.table.Where("id=?", obj.ID).Updates(map[string]interface{}{"name": obj.Name, "telephone": obj.Telephone, "timestamp": obj.Timestamp, "enable": obj.Enable, "wechat_user_id": obj.WechatUserID}).Error
	return errors.WithMessage(err, "dao update Receiver")
}

func (c *Receiver) Pager(pageIndex, pageSize int) (result []*model.Receiver, count int64) {
	res := c.table.Select("*")
	res.Count(&count)
	res.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&result)
	return
}

func (c *Receiver) Delete(id int32) error {
	err := c.table.Delete(&model.Receiver{}, id).Error
	return errors.WithMessage(err, "dao Delete Receiver")
}
func (c *Receiver) Search(name string) (result []*model.Receiver, err error) {
	err = c.table.Where("name like ?", name+"%").Scan(&result).Error
	return
}

func (c *Receiver) GetReceiverByID(id int) (result *model.Receiver, err error) {
	err = c.table.Where("id=?", id).Scan(&result).Error
	return
}

func (c *Receiver) Counter() (count int64) {
	res := c.table.Select("*")
	res.Count(&count)

	return
}
