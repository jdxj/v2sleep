package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestContains(t *testing.T) {
	str1 := "abc香港def"
	str2 := "香港"
	res := strings.Contains(str1, str2)
	fmt.Printf("%t\n", res)
}
