package v2rayng

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type V2rayVmess struct {
	Version   json.RawMessage `json:"v"`
	Name      string          `json:"ps"`
	Address   string          `json:"add"`
	Port      json.RawMessage `json:"port"`
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
