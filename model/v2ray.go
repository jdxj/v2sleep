package model

import (
	"context"
	"time"

	"github.com/jdxj/v2sleep/dao"
)

type AddV2raySubAddrReq struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type AddV2raySubAddrRsp struct {
	ID uint32 `json:"id"`
}

func AddV2raySubAddr(ctx context.Context, req *AddV2raySubAddrReq) (*AddV2raySubAddrRsp, error) {
	now := time.Now()
	sc := &dao.SubConfig{
		Name:     req.Name,
		Type:     uint8(V2raySubAddr),
		Data:     []byte(req.Address),
		CreateAt: now,
		UpdateAt: now,
	}
	err := dao.DB.WithContext(ctx).
		Create(sc).
		Error
	if err != nil {
		return nil, err
	}

	return &AddV2raySubAddrRsp{ID: sc.ID}, nil
}