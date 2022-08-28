package proto

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
)

func NewShareLinkParser() *ShareLinkParser {
	return &ShareLinkParser{}
}

type ShareLinkParser struct {
	v2raies []V2ray
}

func (slp *ShareLinkParser) Encode() ([]byte, error) {
	return encodeV2ray(slp.v2raies)
}

func encodeV2ray(v2raies []V2ray) ([]byte, error) {
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

	slp.v2raies = append(slp.v2raies, v2)
	return nil
}

func (slp *ShareLinkParser) ToV2ray() []V2ray {
	return slp.v2raies
}

func decodeShareLink(data []byte) (V2ray, error) {
	i := strings.Index(string(data), "://")
	if i < 0 {
		return nil, fmt.Errorf("invalid share link: %s", data)
	}

	var (
		v2  V2ray
		err error
	)
	switch string(data[:i]) {
	case "ss":
		v2 = &V2rayShadowsocks{}
		err = v2.Decode(data)
	case "vmess":
		v2 = &V2rayVmess{}
		err = v2.Decode(data)
	case "trojan":
		v2 = &V2rayTrojan{}
		err = v2.Decode(data)
	default:
		return nil, fmt.Errorf("%s not supported: %s", data[:i], data)
	}
	if err != nil {
		return nil, fmt.Errorf("decode share link %s err: %s", data, err)
	}
	return v2, nil
}
