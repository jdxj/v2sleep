package v2rayng

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/jdxj/v2sleep/proto"
	"github.com/jdxj/v2sleep/proto/v2raycore"
)

func NewSubAddrParser() *SubAddrParser {
	return &SubAddrParser{
		hc: &http.Client{},
	}
}

type SubAddrParser struct {
	hc      *http.Client
	V2raies []proto.V2rayNG

	TagPrefix string
}

func (vsa *SubAddrParser) Encode() ([]byte, error) {
	return encodeV2ray(vsa.V2raies)
}

func (vsa *SubAddrParser) Decode(data []byte) error {
	rsp, err := vsa.hc.Get(string(data))
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, rsp.Body)
	if err != nil {
		_ = rsp.Body.Close()
		return err
	}
	_ = rsp.Body.Close()

	noPaddingData := strings.TrimRight(buf.String(), "=")
	data, err = base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(noPaddingData)
	if err != nil {
		return err
	}
	buf.Reset()
	_, err = buf.Write(data)
	if err != nil {
		return err
	}
	for {
		data, err = buf.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		data = bytes.TrimSuffix(data, []byte{'\n'})

		v2, err := decodeShareLink(data)
		if err != nil {
			logrus.Warnf("%s", err)
			continue
		}
		vsa.V2raies = append(vsa.V2raies, v2)
	}
	return nil
}

func (vsa *SubAddrParser) Merge(vs ...proto.V2rayNG) {
	vsa.V2raies = append(vsa.V2raies, vs...)
}

type Filter func(out *v2raycore.Outbound) bool

func (vsa *SubAddrParser) Outbounds(filters ...Filter) ([]byte, error) {
	outs := v2raycore.OutboundConfig{}
	for _, v := range vsa.V2raies {
		c, ok := v.(v2raycore.Outbounder)
		if !ok {
			logrus.Warnf("can not switch to outbounder: %+v", v)
			continue
		}
		out, err := c.Outbound()
		if err != nil {
			return nil, fmt.Errorf("get outbound err: %s", err)
		}

		add := true
		for _, f := range filters {
			if !f(out) {
				add = false
				break
			}
		}

		if add {
			outs.Outbounds = append(outs.Outbounds, out)
		}
	}
	return json.MarshalIndent(outs, "", "  ")
}
