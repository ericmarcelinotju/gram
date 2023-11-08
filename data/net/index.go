package net

import (
	"crypto/tls"
	"time"

	"net/http"
	"net/url"

	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/library/net"
)

func Init() (*net.NetClient, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Timeout:   time.Duration(config.Get().Net.Timeout),
		Transport: tr,
	}

	var httpUrl *url.URL
	return net.New(httpUrl.String(), "username", "password", client)
}
