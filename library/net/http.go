package net

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type NetClient struct {
	HttpURL  string
	WsURL    string
	Client   http.Client
	Username string
	Password string
	Token    *string
	quit     chan string
	IsOnline bool
}

func New(httpUrl, username, password string, client http.Client) (*NetClient, error) {
	netClient := &NetClient{
		HttpURL:  httpUrl,
		Client:   client,
		Username: username,
		Password: password,
		quit:     make(chan string),
	}

	//Login to get token
	err := netClient.login()
	if err == nil {
		netClient.IsOnline = true
	}

	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				//Heartbeat to update token
				if netClient.IsOnline {
					err := netClient.heartbeat()
					if err != nil {
						netClient.IsOnline = false
					} else {
						netClient.IsOnline = true
					}
				}

				if !netClient.IsOnline {
					err := netClient.login()
					if err != nil {
						netClient.IsOnline = false
					} else {
						netClient.IsOnline = true
					}
				}

			case <-netClient.quit:
				ticker.Stop()
				return
			}
		}
	}()

	return netClient, nil
}

func (h *NetClient) setHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	if h.Token != nil {
		req.Header.Set("Authorization", "Bearer "+*h.Token)
	}
}

func (h *NetClient) Get(path string, data map[string]interface{}) (body []byte, err error) {
	urlValues := url.Values{}

	for key, value := range data {
		urlValues.Set(key, fmt.Sprintf("%v", value))
	}

	req, err := http.NewRequest("GET", h.HttpURL+path+"?"+urlValues.Encode(), nil)
	if err != nil {
		return body, err
	}
	h.setHeaders(req)
	// Send request
	resp, err := h.Client.Do(req)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	if resp.StatusCode != 200 {
		err = errors.New(string(body))
	}

	return body, err
}

func (h *NetClient) Post(path string, data map[string]interface{}) (body []byte, err error) {
	jsonString, err := json.Marshal(data)
	if err != nil {
		return body, err
	}

	payload := strings.NewReader(string(jsonString))
	if err != nil {
		return body, err
	}

	req, err := http.NewRequest("POST", h.HttpURL+path, payload)
	if err != nil {
		return body, err
	}
	h.setHeaders(req)
	// Send request
	resp, err := h.Client.Do(req)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	if resp.StatusCode != 200 {
		err = errors.New(string(body))
	}

	return body, err
}

func (h *NetClient) Put(path string, data map[string]interface{}) (body []byte, err error) {
	jsonString, err := json.Marshal(data)
	if err != nil {
		return body, err
	}

	payload := strings.NewReader(string(jsonString))

	req, err := http.NewRequest("PUT", h.HttpURL+path, payload)
	if err != nil {
		return body, err
	}
	h.setHeaders(req)
	// Send request
	resp, err := h.Client.Do(req)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	if resp.StatusCode != 200 {
		err = errors.New(string(body))
	}

	return body, err
}

func (h *NetClient) Patch(path string, data map[string]interface{}) (body []byte, err error) {
	jsonString, err := json.Marshal(data)
	if err != nil {
		return body, err
	}

	payload := strings.NewReader(string(jsonString))

	req, err := http.NewRequest("PATCH", h.HttpURL+path, payload)
	if err != nil {
		return body, err
	}
	h.setHeaders(req)
	// Send request
	resp, err := h.Client.Do(req)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	if resp.StatusCode != 200 {
		err = errors.New(string(body))
	}

	return body, err
}
