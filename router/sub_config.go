package router

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/jdxj/v2sleep/model"
)

func ListSubConfig(ctx *gin.Context) {
	handle(ctx, nil, func(ctx context.Context) (interface{}, error) {
		return model.ListSubConfig(ctx)
	})
}

func DeleteSubConfig(ctx *gin.Context) {
	req := &model.DeleteSubConfigReq{}
	handle(ctx, req, func(ctx context.Context) (interface{}, error) {
		return nil, model.DeleteSubConfig(ctx, req)
	})
}

func AddClashSubAddr(ctx *gin.Context) {
	req := &model.AddSubAddrReq{}
	handle(ctx, req, func(ctx context.Context) (interface{}, error) {
		return model.AddClashSubAddr(ctx, req)
	})
}

func AddV2raySubAddr(ctx *gin.Context) {
	req := &model.AddSubAddrReq{}
	handle(ctx, req, func(ctx context.Context) (interface{}, error) {
		return model.AddV2raySubAddr(ctx, req)
	})
}
