package model

import (
	"bytes"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jdxj/v2sleep/dao"
	"github.com/jdxj/v2sleep/proto"
)

func GenShare(ctx context.Context) (*bytes.Buffer, error) {
	var scs []*dao.SubConfig
	err := dao.DB.WithContext(ctx).
		Model(dao.SubConfig{}).
		Find(&scs).
		Error
	if err != nil {
		return nil, err
	}

	vsa := proto.NewV2raySubAddrParser()
	for _, sc := range scs {
		switch proto.ConfType(sc.Type) {
		case proto.ClashSubAddr:
			logrus.Warnf("clash sub addr not support")
			continue
		case proto.V2raySubAddr:
			err = vsa.Decode(sc.Data)
		}
		if err != nil {
			return nil, fmt.Errorf("name: %s, err: %w", sc.Name, err)
		}
	}

	data, err := vsa.Encode()
	return bytes.NewBuffer(data), err
}
