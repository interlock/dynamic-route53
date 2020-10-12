package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func lookup(network string) (string, error) {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _network, addr string) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	}
	httpClient := &http.Client{
		Transport: transport,
	}
	res, err := httpClient.Get("https://ifconfig.co/ip")
	if err != nil {
		log.Fatal(err)
	}
	ip, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return string(ip), err
}
