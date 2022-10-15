package v2rayng

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/jdxj/v2sleep/proto"
)

func NewShareLinkParser() *ShareLinkParser {
	return &ShareLinkParser{}
}

type ShareLinkParser struct {
	V2raies []proto.V2rayNG
}

func (slp *ShareLinkParser) Encode() ([]byte, error) {
	return encodeV2ray(slp.V2raies)
}

func encodeV2ray(v2raies []proto.V2rayNG) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	for _, v := range v2raies {
		data, err := v.Encode()
		if err != nil {
			return nil, err
		}
		data = append(data, '\n')
		_, err = buf.Write(data)
		if err != nil {
			return nil, err
		}
	}

	s := base64.StdEncoding.EncodeToString(buf.Bytes())
	buf.Reset()
	_, err := buf.WriteString(s)
	return buf.Bytes(), err
}

func (slp *ShareLinkParser) Decode(data []byte) error {
	v2, err := decodeShareLink(data)
	if err != nil {
		return err
	}

	slp.V2raies = append(slp.V2raies, v2)
	return nil
}

func decodeShareLink(data []byte) (proto.V2rayNG, error) {
	i := strings.Index(string(data), "://")
	if i < 0 {
		return nil, fmt.Errorf("invalid share link: %s", data)
	}

	var (
		v2  proto.V2rayNG
		err error
	)
	switch string(data[:i]) {
	case "ss":
		v2 = &Shadowsocks{}
		err = v2.Decode(data)
	case "vmess":
		v2 = &Vmess{}
		err = v2.Decode(data)
	case "trojan":
		v2 = &Trojan{}
		err = v2.Decode(data)
	default:
		return nil, fmt.Errorf("%s not supported: %s", data[:i], data)
	}
	if err != nil {
		return nil, fmt.Errorf("decode share link %s err: %s", data, err)
	}
	return v2, nil
}
