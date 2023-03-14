package util

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

func EvalHttpExist(url string) bool {
	// 1. 定义一个Client对象，并自定义传输，以设置超时处理。
	isexist := false
	c := http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				timeout := time.Second * 2
				return net.DialTimeout(network, addr, timeout)
			},
		},
	}

	// 2.  发送Head请求。
	resp, err := c.Head(url)
	if err != nil {
		fmt.Printf("head %s failed, err:%v\n", url, err)
	} else {
		fmt.Printf("%s head success, status:%v\n", url, resp.Status)
		isexist = true
	}
	return isexist
}
