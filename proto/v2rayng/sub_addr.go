package v2rayng

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/jdxj/v2sleep/proto"
)

func NewSubAddrParser() *SubAddrParser {
	return &SubAddrParser{
		hc: &http.Client{},
	}
}

type SubAddrParser struct {
	hc      *http.Client
	V2raies []proto.V2rayNG
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

	data, err = base64.StdEncoding.DecodeString(buf.String())
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
