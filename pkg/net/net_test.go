package net

import (
	"fmt"
	"testing"
)

func TestIpToInt(t *testing.T) {

	a := IpToInt("192.168.0.1")
	fmt.Println(a)
}
