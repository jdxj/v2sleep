package model

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"time"

	"github.com/jdxj/v2sleep/dao"
)

type AddSigReq struct {
	Sig string `json:"sig"`
}

func AddSig(ctx context.Context, req *AddSigReq) error {
	req.Sig = strings.TrimSpace(req.Sig)
	if req.Sig == "" {
		h := md5.New()
		h.Write([]byte(time.Now().String()))
		req.Sig = hex.EncodeToString(h.Sum(nil))
	}
	now := time.Now()
	return dao.DB.WithContext(ctx).
		Create(&dao.Section{
			Sig:      req.Sig,
			CreateAt: now,
			UpdateAt: now,
		}).
		Error
}

type GetSigRsp struct {
	Sig string `json:"sig"`
}

func GetSig(ctx context.Context) (*GetSigRsp, error) {
	var sig string
	err := dao.DB.WithContext(ctx).
		Select("sig").
		Model(dao.Section{}).
		Order("update_at DESC").
		First(&sig).
		Error
	if err != nil {
		return nil, err
	}
	return &GetSigRsp{Sig: sig}, nil
}
