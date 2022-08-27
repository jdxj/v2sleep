package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jdxj/v2sleep/dao"
	"github.com/jdxj/v2sleep/proto"
)

type item struct {
	ID   uint32         `json:"id"`
	Name string         `json:"name"`
	Type proto.ConfType `json:"type"`
	Data string         `json:"data"`
}

type ListSubConfigRsp struct {
	Count int    `json:"count"`
	List  []item `json:"list"`
}

func ListSubConfig(ctx context.Context) (*ListSubConfigRsp, error) {
	var scs []dao.SubConfig
	err := dao.DB.WithContext(ctx).
		Model(dao.SubConfig{}).
		Find(&scs).
		Error
	if err != nil {
		return nil, err
	}
	rsp := &ListSubConfigRsp{
		Count: len(scs),
	}
	for _, v := range scs {
		rsp.List = append(rsp.List, item{
			ID:   v.ID,
			Name: v.Name,
			Type: proto.ConfType(v.Type),
			Data: string(v.Data),
		})
	}
	return rsp, nil
}

type AddClashSubAddrReq struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type AddClashSubAddrRsp struct {
	ID uint32 `json:"id"`
}

func AddClashSubAddr(ctx context.Context, req *AddClashSubAddrReq) (*AddClashSubAddrRsp, error) {
	now := time.Now()
	sc := &dao.SubConfig{
		Name:     req.Name,
		Type:     uint8(proto.ClashSubAddr),
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
	return &AddClashSubAddrRsp{ID: sc.ID}, nil
}

type DeleteSubConfigReq struct {
	ID uint32 `json:"id" binding:"required"`
}

func DeleteSubConfig(ctx context.Context, req *DeleteSubConfigReq) error {
	result := dao.DB.WithContext(ctx).
		Delete(dao.SubConfig{}, req.ID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("id %d not found", req.ID)
	}
	return nil
}
