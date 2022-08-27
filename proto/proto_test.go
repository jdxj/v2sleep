package proto

import (
	"fmt"
	"net/url"
	"regexp"
	"testing"
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
