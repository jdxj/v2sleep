package dao

import "time"

type SubConfig struct {
	ID       uint32
	Name     string
	Type     uint8
	Data     []byte
	CreateAt time.Time
	UpdateAt time.Time
}

func (SubConfig) TableName() string {
	return "sub_config"
}
