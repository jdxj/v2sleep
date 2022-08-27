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
	buf := bytes.NewBuffer(nil)
	for _, v := range vsa.v2raies {
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

		i := strings.Index(string(data), "://")
		if i < 0 {
			logrus.Warnf("read invalid row: %s", data)
			continue
		}

		var v2 V2ray
		switch string(data[:i]) {
		case "ss":
			v2 = &V2rayShadowsocks{}
			err = v2.Decode(data)
		case "vmess":
			v2 = &V2rayVmess{}
			err = v2.Decode(data)
		default:
			logrus.Warnf("%s not supported: %s", data[:i], data)
			continue
		}
		if err != nil {
			return err
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
