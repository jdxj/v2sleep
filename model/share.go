package model

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/jdxj/v2sleep/dao"
)

func GenShare(ctx context.Context) (*bytes.Buffer, error) {
	// todo: 暂时使用一个
	sc := &dao.SubConfig{}
	err := dao.DB.WithContext(ctx).
		Model(dao.SubConfig{}).
		// todo: 暂时使用 clash
		Where("type = ?", V2raySubAddr).
		First(&sc).
		Error
	if err != nil {
		return nil, err
	}

	rsp, err := http.Get(string(sc.Data))
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, rsp.Body)
	return buf, err
}
