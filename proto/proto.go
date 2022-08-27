package proto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"gopkg.in/yaml.v3"
)

var (
	ParserMap = map[ConfType]Parser{
		V2raySubAddr: &V2raySubAddrParser{hc: &http.Client{}},
	}
)

type Parser interface {
	Parse(data []byte) (string, error)
}

type V2raySubAddrParser struct {
	hc *http.Client
}

func (vsa *V2raySubAddrParser) Parse(data []byte) (string, error) {
	rsp, err := vsa.hc.Get(string(data))
	if err != nil {
		return "", err
	}
	_ = rsp
	// todo
	return "", nil

}

func NewClashConfig(r io.Reader) (*ClashConfig, error) {
	decoder := yaml.NewDecoder(r)
	cc := &ClashConfig{}
	return cc, decoder.Decode(cc)
}

type Proxy struct {
	Name     string `yaml:"name"`
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	Type     string `yaml:"type"`
	Cipher   string `json:"cipher"`
	Password string `json:"password"`
}

type ClashConfig struct {
	Proxies []Proxy `yaml:"proxies"`
}

func (cc *ClashConfig) ToV2rayShadowsocks() []*V2rayShadowsocks {
	var vsss []*V2rayShadowsocks
	for _, p := range cc.Proxies {
		vsss = append(vsss, &V2rayShadowsocks{
			Cipher:   p.Cipher,
			Password: p.Password,
			Server:   p.Server,
			Port:     p.Port,
			Name:     p.Name,
		})
	}
	return vsss
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
	Version   string          `json:"v"`
	Name      string          `json:"ps"`
	Address   string          `json:"add"`
	Port      string          `json:"port"`
	ID        string          `json:"id"`
	AID       json.RawMessage `json:"aid"`
	Security  string          `json:"scy"`
	TransType string          `json:"net"`
	FakeType  string          `json:"type"`
	FakeHost  string          `json:"host"`
	FakePath  string          `json:"path"`
	TLS       string          `json:"tls"`
	SNI       string          `json:"sni"`
}

func (vv *V2rayVmess) String() string {
	// todo
	return ""
}
