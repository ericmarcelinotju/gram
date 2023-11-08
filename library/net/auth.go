package net

import (
	"encoding/json"
	"errors"
	"fmt"
)

type LoginResponse struct {
	Status string `json:"status"`
	Err    string `json:"errno"`
	Token  string `json:"token"`
}

func (n *NetClient) login() (err error) {
	body, err := n.Post("/login", map[string]interface{}{
		"username": n.Username,
		"password": n.Password,
	})
	if err != nil {
		return err
	}

	var loginResp LoginResponse
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		return err
	}

	n.Token = &loginResp.Token

	return nil
}

func (n *NetClient) heartbeat() (err error) {
	if n.Token == nil {
		return errors.New("no token")
	}
	body, err := n.Post(fmt.Sprintf("/heartbeat?token=%s", *n.Token), nil)
	if err != nil {
		return err
	}

	var loginResp LoginResponse
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		return err
	}

	return nil
}
