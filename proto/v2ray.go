package proto

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func NewV2raySubAddrParser() *V2raySubAddrParser {
	return &V2raySubAddrParser{
		hc: &http.Client{},
	}
}

type V2raySubAddrParser struct {
	hc      *http.Client
	v2raies []V2ray
}

func (vsa *V2raySubAddrParser) Encode() ([]byte, error) {
	return encodeV2ray(vsa.v2raies)
}

func (vsa *V2raySubAddrParser) Decode(data []byte) error {
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
		vsa.v2raies = append(vsa.v2raies, v2)
	}
	return nil
}

func (vsa *V2raySubAddrParser) Merge(vs ...V2ray) {
	vsa.v2raies = append(vsa.v2raies, vs...)
}

type V2ray interface {
	Encode() ([]byte, error)
	Decode([]byte) error
}

type V2rayShadowsocks struct {
	Cipher   string
	Password string
	Server   string
	Port     int
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

	data, err = base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(u.Host)
	if err != nil {
		return err
	}
	reg := regexp.MustCompile(`(.*):(.*)@(.*):(.*)`)
	conf := reg.FindStringSubmatch(string(data))
	if len(conf) != 5 {
		return fmt.Errorf("invalid ss config: %s", data)
	}

	vss.Cipher = conf[1]
	vss.Password = conf[2]
	vss.Server = conf[3]
	port, err := strconv.ParseInt(conf[4], 10, 64)
	if err != nil {
		return fmt.Errorf("parse port %s err: %s", conf[4], err)
	}
	vss.Port = int(port)
	return nil
}

type V2rayVmess struct {
	Version   string `json:"v"`
	Name      string `json:"ps"`
	Address   string `json:"add"`
	Port      int    `json:"port"`
	ID        string `json:"id"`
	AID       int    `json:"aid"`
	Security  string `json:"scy"`
	TransType string `json:"net"`
	FakeType  string `json:"type"`
	FakeHost  string `json:"host"`
	FakePath  string `json:"path"`
	TLS       string `json:"tls"`
	SNI       string `json:"sni"`
}

func (vv *V2rayVmess) Encode() ([]byte, error) {
	data, err := json.Marshal(vv)
	if err != nil {
		return nil, err
	}
	s := base64.StdEncoding.EncodeToString(data)
	return []byte(fmt.Sprintf("vmess://%s", s)), nil
}

func (vv *V2rayVmess) Decode(data []byte) error {
	s := strings.TrimPrefix(string(data), "vmess://")
	jsonData, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, vv)
}

type V2rayTrojan struct {
	// path
	Password string
	Server   string
	Port     int
	// query
	Security   string
	HeaderType string
	Type       string
	// fragment
	Name string
}

func (vt *V2rayTrojan) Encode() ([]byte, error) {
	return []byte(fmt.Sprintf("trojan://%s@%s:%d?security=%s&headerType=%s&type=%s#%s",
		vt.Password, vt.Server, vt.Port, vt.Security, vt.HeaderType, vt.Type, url.PathEscape(vt.Name),
	)), nil
}

func (vt *V2rayTrojan) Decode(data []byte) error {
	u, err := url.Parse(string(data))
	if err != nil {
		return err
	}
	vt.Name, err = url.PathUnescape(u.Fragment)
	if err != nil {
		return err
	}

	vt.Password = u.User.Username()
	vt.Server = u.Hostname()
	port, err := strconv.ParseInt(u.Port(), 10, 64)
	if err != nil {
		return err
	}
	vt.Port = int(port)

	vt.Security = u.Query().Get("security")
	vt.HeaderType = u.Query().Get("headerType")
	vt.Type = u.Query().Get("type")
	return nil
}
