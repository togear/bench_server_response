package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

type QueryResp struct {
	Addr string
	Time float64
}

func query(ip string, port string, c chan QueryResp) {
	start_ts := time.Now()
	var timeout = time.Duration(15 * time.Second)
	host := fmt.Sprintf("%s:%s", ip, port)
//	url_proxy := &url.URL{Host: host}
	http_proxy_url := fmt.Sprintf("http://%s:%s/",ip,port)

	url_i := url.URL{}
	url_proxy, _ := url_i.Parse(http_proxy_url)
//	fmt.Println(http_proxy_url)

	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(url_proxy)},
		Timeout:   timeout}

	resp, err := client.Get("http://www.test.com/benchmark/test.m3u8")
	if err != nil {
		c <- QueryResp{Addr: host, Time: float64(-1)}
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	time_diff := time.Now().UnixNano() - start_ts.UnixNano()
//	判断成功标志，这里是一个例子，取得的是ts文件
	if strings.Contains(string(body), "ts") {
		c <- QueryResp{Addr: host, Time: float64(time_diff) / 1e9}
	} else {
		c <- QueryResp{Addr: host, Time: float64(-1)}
	}
}

func main() {
	dat, _ := ioutil.ReadFile("myip.lst")
//	dat, _ := ioutil.ReadFile("ip.txt.all")
	dats := strings.Split(strings.TrimSuffix(string(dat), "\n"), "\n")

	runtime.GOMAXPROCS(4)

	resp_chan := make(chan QueryResp, 170)

	for _, addr := range dats {
		addrs := strings.SplitN(addr, string(' '), 2)
		ip, port := addrs[0], addrs[1]
		go query(ip, port, resp_chan)
	}

	for _, _ = range dats {
		r := <-resp_chan
//		if r.Time > 1e-9 {
			fmt.Printf("%s %v\n", r.Addr, r.Time)
//		}
	}
}
