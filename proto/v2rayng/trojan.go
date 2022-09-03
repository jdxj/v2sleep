package v2rayng

import (
	"fmt"
	"net/url"
	"strconv"
)

type V2rayTrojan struct {
	// path
	Password string
	Server   string
	Port     int64
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
	vt.Port, err = strconv.ParseInt(u.Port(), 10, 64)
	if err != nil {
		return err
	}

	vt.Security = u.Query().Get("security")
	vt.HeaderType = u.Query().Get("headerType")
	vt.Type = u.Query().Get("type")
	return nil
}
