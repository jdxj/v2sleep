package model

import (
	"bytes"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jdxj/v2sleep/dao"
	"github.com/jdxj/v2sleep/proto"
	"github.com/jdxj/v2sleep/proto/clash"
	"github.com/jdxj/v2sleep/proto/convert"
	"github.com/jdxj/v2sleep/proto/v2rayng"
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

	csa := clash.NewSubAddrParser()
	slp := v2rayng.NewShareLinkParser()
	vsa := v2rayng.NewSubAddrParser()
	for _, sc := range scs {
		switch proto.ConfType(sc.Type) {
		case proto.ClashSubAddr:
			err = csa.Decode(sc.Data)
		case proto.ShareLink:
			err = slp.Decode(sc.Data)
		case proto.V2raySubAddr:
			err = vsa.Decode(sc.Data)
		default:
			logrus.Warnf("conf type %d not support", sc.Type)
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("name: %s, err: %w", sc.Name, err)
		}
	}

	vsa.Merge(slp.V2raies...)
	vsa.Merge(convert.ProxyToV2rayNG(csa.Proxies...)...)
	data, err := vsa.Encode()
	return bytes.NewBuffer(data), err
}
