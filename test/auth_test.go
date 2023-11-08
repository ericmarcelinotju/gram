package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	"gitlab.com/firelogik/helios/router/http/auth"
	dto "gitlab.com/firelogik/helios/router/http/dto/responses"

	response "gitlab.com/firelogik/helios/utils/http"
)

var AUTH_TOKEN = ""

func TestLoginHandler(t *testing.T) {
	r := SetUpRouter()
	svc := AUTHSVC

	r.POST("/api/auth/login", auth.Login(svc))
	request, _ := json.Marshal(auth.LoginValidator{
		Username: "ninachristian",
		Password: "qwerty123",
	})
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(request))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := &dto.LoginResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	AUTH_TOKEN = resp.Token

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLogoutHandler(t *testing.T) {
	r := SetUpRouter()
	svc := AUTHSVC

	r.POST("/api/auth/logout", auth.Logout(svc))
	req, _ := http.NewRequest("POST", "/api/auth/logout", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authorization", AUTH_TOKEN)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp response.SetResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, http.StatusOK, w.Code)
}
