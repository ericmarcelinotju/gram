package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ericmarcelinotju/gram/domain/model"
	"github.com/ericmarcelinotju/gram/router/http/dto"
	responses "github.com/ericmarcelinotju/gram/router/http/dto/responses"
	"github.com/ericmarcelinotju/gram/router/http/user"
	"github.com/go-playground/assert/v2"

	response "github.com/ericmarcelinotju/gram/utils/http"
)

func TestListUserHandler(t *testing.T) {
	r := SetUpRouter()
	svc := USERSVC

	r.GET("/api/user", user.GetUser(svc))
	req, _ := http.NewRequest("GET", "/api/user", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := &responses.ListUserResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, 0, len(resp.Users))
}

func TestListPaginationUserHandler(t *testing.T) {
	r := SetUpRouter()
	svc := USERSVC

	r.GET("/api/user", user.GetUser(svc))

	currentPage := 1
	totalPerPage := 10

	request, _ := json.Marshal(user.UserFilter{
		Pagination: dto.Pagination{
			CurrentPage:  &currentPage,
			TotalPerPage: &totalPerPage,
		},
	})
	req, _ := http.NewRequest("GET", "/api/user", bytes.NewBuffer(request))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := &responses.ListUserResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, totalPerPage, len(resp.Users))
}

func TestDetailUserHandler(t *testing.T) {
	r := SetUpRouter()
	svc := USERSVC

	users, _, _ := USERREPO.SelectUser(context.Background(), nil)
	currUser := users[0]

	r.GET("/api/user/:id", user.GetUserDetail(svc))

	req, _ := http.NewRequest("GET", "/api/user/"+currUser.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := &responses.UserResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, currUser.ID, resp.ID)
}

func TestCreateUserHandler(t *testing.T) {
	r := SetUpRouter()
	svc := USERSVC

	roles, _, _ := ROLEREPO.SelectRole(context.Background(), nil)
	currRole := roles[0]

	r.POST("/api/user", user.PostUser(svc))

	userName := "testuser"
	payload := user.PostUserPayload{
		Username: userName,
		Email:    "testing@gmail.com",
		Password: "password",
		RoleID:   currRole.ID,
	}
	payloadJson, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/user", bytes.NewBuffer(payloadJson))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := &responses.UserResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, userName, resp.Username)
}

func TestUpdateUserHandler(t *testing.T) {
	r := SetUpRouter()
	svc := USERSVC

	users, _, _ := USERREPO.SelectUser(context.Background(), &model.User{Username: "testuser"})
	currUser := users[0]

	r.PUT("/api/user/:id", user.PutUser(svc))

	userName := "testchangeuser"
	payload := user.PutUserPayload{
		Username: userName,
	}
	payloadJson, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", "/api/user/"+currUser.ID, bytes.NewBuffer(payloadJson))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := &responses.UserResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, userName, resp.Username)
}

func TestDeleteUserHandler(t *testing.T) {
	r := SetUpRouter()
	svc := USERSVC

	users, _, _ := USERREPO.SelectUser(context.Background(), &model.User{Username: "testchangeuser"})
	currUser := users[0]

	r.DELETE("/api/user/:id", user.DeleteUser(svc))

	req, _ := http.NewRequest("DELETE", "/api/user/"+currUser.ID, nil)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := ""
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)

	afterDelUsers, _, _ := USERREPO.SelectUser(context.Background(), &model.User{Username: "testchangeuser"})
	assert.Equal(t, 0, len(afterDelUsers))
}
