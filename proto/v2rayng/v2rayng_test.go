package v2rayng

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/jdxj/v2sleep/proto/clash"
)

func TestV2rayShadowsocks_String(t *testing.T) {
}

func TestParseURL(t *testing.T) {
	u, err := url.Parse("ss://abc#def")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("path: %s\n", u.Path)
	fmt.Printf("raw path: %s\n", u.RawPath)
	fmt.Printf("fragment: %s\n", u.Fragment)
	fmt.Printf("raw fragment: %s\n", u.RawFragment)
	fmt.Printf("host: %s\n", u.Host)
	fmt.Printf("raw path: %s\n", u.RawPath)
	fmt.Printf("scheme: %s\n", u.Scheme)
}

func TestScan(t *testing.T) {
	r := "abc:def@123:456"
	reg, err := regexp.Compile(`(.*):(.*)@(.*):(.*)`)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	ss := reg.FindStringSubmatch(r)
	fmt.Printf("%+v\n", ss)
}

func TestV2rayShadowsocks_Decode(t *testing.T) {
	vss := &V2rayShadowsocks{}
	data := ""
	err := vss.Decode([]byte(data))
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("%+v\n", vss)
	data2, err := vss.Encode()
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	if data == string(data2) {
		fmt.Printf("ok\n")
	}
}

func TestV2rayVmess_Encode(t *testing.T) {
	vv := &V2rayVmess{}
	data := ""
	err := vv.Decode([]byte(data))
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("%+v\n", vv)

	data2, err := vv.Encode()
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	if data == string(data2) {
		fmt.Printf("ok\n")
	} else {
		fmt.Printf("%s\n", data2)
	}
}

func TestNewV2raySubAddrParser(t *testing.T) {
	vsa := NewSubAddrParser()
	addr := ""
	err := vsa.Decode([]byte(addr))
	if err != nil {
		t.Fatalf("%s\n", err)
	}
}

func TestNewClashSubAddrParser(t *testing.T) {
	csa := clash.NewSubAddrParser()
	addr := ""
	err := csa.Decode([]byte(addr))
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	for _, v := range csa.Proxies {
		fmt.Printf("%+v\n", v)
	}
}

func TestV2rayTrojan_Decode(t *testing.T) {
	vt := &V2rayTrojan{}
	data := ""
	err := vt.Decode([]byte(data))
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("%+v\n", vt)

	data2, err := vt.Encode()
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	if data == string(data2) {
		fmt.Printf("ok\n")
	}
}

func TestTrim(t *testing.T) {
	ts := []struct {
		name string
		cas  string
		want string
	}{
		{
			name: "case1",
			cas:  "abc==",
			want: "abc",
		},
		{
			name: "case2",
			cas:  "def=",
			want: "def",
		},
		{
			name: "case3",
			cas:  "ghi",
			want: "ghi",
		},
	}

	for _, v := range ts {
		t.Run(v.name, func(t *testing.T) {
			if s := strings.TrimRight(v.cas, "="); s != v.want {
				t.Fatalf("name: %s, get: %s, want: %s",
					v.name, s, v.want)
			}
		})
	}
}

func TestNewSubAddrParser(t *testing.T) {
	sap := NewSubAddrParser()
	data := ""
	err := sap.Decode([]byte(data))
	if err != nil {
		t.Fatalf("%s\n", err)
	}
}
