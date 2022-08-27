package dao

import (
	"time"
)

type Section struct {
	ID       uint32
	Sig      string
	CreateAt time.Time
	UpdateAt time.Time
}

func (Section) TableName() string {
	return "section"
}
