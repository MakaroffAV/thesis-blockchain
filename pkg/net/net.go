package net

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func IpToInt(host string) uint32 {

	ip := strings.Split(host, ":")[0]

	parsedIp := net.ParseIP(ip)

	a := parsedIp.To4()

	return uint32(a[0])<<24 | uint32(a[1])<<16 | uint32(a[2])<<8 | uint32(a[3])

}

func GetIp() string {

	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	return string(ip)

}
