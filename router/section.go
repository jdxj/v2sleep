package router

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/jdxj/v2sleep/model"
)

func SetSig(ctx *gin.Context) {
	req := &model.AddSigReq{}
	handle(ctx, req, func(ctx context.Context) (interface{}, error) {
		return nil, model.AddSig(ctx, req)
	})
}

func GetSig(ctx *gin.Context) {
	handle(ctx, nil, func(ctx context.Context) (interface{}, error) {
		return model.GetSig(ctx)
	})
}
