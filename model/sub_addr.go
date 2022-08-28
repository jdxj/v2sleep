package model

import (
	"context"
	"time"

	"github.com/jdxj/v2sleep/dao"
	"github.com/jdxj/v2sleep/proto"
)

type AddSubAddrReq struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type AddSubAddrRsp struct {
	ID uint32 `json:"id"`
}

func AddShareLink(ctx context.Context, req *AddSubAddrReq) (*AddSubAddrRsp, error) {
	return addSubAddr(ctx, proto.ShareLink, req)
}

func AddClashSubAddr(ctx context.Context, req *AddSubAddrReq) (*AddSubAddrRsp, error) {
	return addSubAddr(ctx, proto.ClashSubAddr, req)
}

func AddV2raySubAddr(ctx context.Context, req *AddSubAddrReq) (*AddSubAddrRsp, error) {
	return addSubAddr(ctx, proto.V2raySubAddr, req)
}

func addSubAddr(ctx context.Context, ct proto.ConfType, req *AddSubAddrReq) (*AddSubAddrRsp, error) {
	now := time.Now()
	sc := &dao.SubConfig{
		Name:     req.Name,
		Type:     uint8(ct),
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
	return &AddSubAddrRsp{ID: sc.ID}, nil
}
