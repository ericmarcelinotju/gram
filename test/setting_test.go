package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	dto "github.com/ericmarcelinotju/gram/router/http/dto/responses"
	"github.com/ericmarcelinotju/gram/router/http/setting"
	"github.com/go-playground/assert/v2"

	response "github.com/ericmarcelinotju/gram/utils/http"
)

func TestListSettingHandler(t *testing.T) {
	r := SetUpRouter()
	svc := SETTINGSVC

	r.GET("/api/setting", setting.GetSetting(svc))
	req, _ := http.NewRequest("GET", "/api/setting", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := dto.ListSettingResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, 0, len(resp.Settings))
}

func TestCreateSettingHandler(t *testing.T) {
	r := SetUpRouter()
	svc := SETTINGSVC

	r.POST("/api/setting", setting.SaveSetting(svc))

	settingName := "testsetting"
	settingValue := "testsettingvalue"
	payload := setting.SettingPayload{
		Name:  settingName,
		Value: settingValue,
	}
	payloadJson, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/setting", bytes.NewBuffer(payloadJson))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := dto.SettingResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, settingName, resp.Name)
	assert.Equal(t, settingValue, resp.Value)
}

func TestUpdateSettingHandler(t *testing.T) {
	r := SetUpRouter()
	svc := SETTINGSVC

	r.POST("/api/setting", setting.SaveSetting(svc))

	settingName := "testchangesetting"
	settingValue := "testchangesettingvalue"
	payload := setting.SettingPayload{
		Name:  settingName,
		Value: settingValue,
	}
	payloadJson, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/setting", bytes.NewBuffer(payloadJson))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := dto.SettingResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, settingName, resp.Name)
	assert.Equal(t, settingValue, resp.Value)
}
