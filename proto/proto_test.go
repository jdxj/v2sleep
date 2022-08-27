package proto

import (
	"fmt"
	"testing"
)

func TestV2rayShadowsocks_String(t *testing.T) {
	vss := &V2rayShadowsocks{}
	fmt.Printf("%s\n", vss)
}
