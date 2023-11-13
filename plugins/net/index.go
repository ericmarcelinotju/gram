package net

import (
	"crypto/tls"
	"time"

	"net/http"
	"net/url"

	"github.com/ericmarcelinotju/gram/config"
)

func NewNetClient() (*NetClient, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Timeout:   time.Duration(config.Get().Net.Timeout),
		Transport: tr,
	}

	var httpUrl *url.URL
	return New(httpUrl.String(), "username", "password", client)
}
