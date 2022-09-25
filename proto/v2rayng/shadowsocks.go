package v2rayng

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type V2rayShadowsocks struct {
	Cipher   string
	Password string
	Server   string
	Port     int64
	Name     string
}

func (vss *V2rayShadowsocks) Encode() ([]byte, error) {
	path := fmt.Sprintf("%s:%s@%s:%d", vss.Cipher, vss.Password, vss.Server, vss.Port)
	path = base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(path))
	anchor := url.PathEscape(vss.Name)
	return []byte(fmt.Sprintf("ss://%s#%s", path, anchor)), nil
}

func (vss *V2rayShadowsocks) Decode(data []byte) error {
	u, err := url.Parse(string(data))
	if err != nil {
		return err
	}
	vss.Name, err = url.PathUnescape(u.Fragment)
	if err != nil {
		return err
	}

	host := strings.TrimRight(u.Host, "=")
	data, err = base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(host)
	if err != nil {
		return err
	}
	u, err = url.Parse(fmt.Sprintf("%s://%s", u.Scheme, data))
	if err != nil {
		return err
	}
	vss.Cipher = u.User.Username()
	vss.Password, _ = u.User.Password()
	vss.Server = u.Hostname()
	vss.Port, err = strconv.ParseInt(u.Port(), 10, 64)
	if err != nil {
		return fmt.Errorf("parse port %s err: %s", u.Port(), err)
	}
	return nil
}
