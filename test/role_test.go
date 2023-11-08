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
	"github.com/ericmarcelinotju/gram/router/http/role"
	"github.com/go-playground/assert/v2"

	response "github.com/ericmarcelinotju/gram/utils/http"
)

func TestListRoleHandler(t *testing.T) {
	r := SetUpRouter()
	svc := ROLESVC

	r.GET("/api/role", role.GetRole(svc))
	req, _ := http.NewRequest("GET", "/api/role", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := &responses.ListRoleResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, 0, len(resp.Roles))
}

func TestListPaginationRoleHandler(t *testing.T) {
	r := SetUpRouter()
	svc := ROLESVC

	r.GET("/api/role", role.GetRole(svc))

	currentPage := 1
	totalPerPage := 10

	request, _ := json.Marshal(role.RoleFilter{
		Pagination: dto.Pagination{
			CurrentPage:  &currentPage,
			TotalPerPage: &totalPerPage,
		},
	})
	req, _ := http.NewRequest("GET", "/api/role", bytes.NewBuffer(request))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := &responses.ListRoleResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, totalPerPage, len(resp.Roles))
}

func TestDetailRoleHandler(t *testing.T) {
	r := SetUpRouter()
	svc := ROLESVC

	roles, _, _ := ROLEREPO.SelectRole(context.Background(), nil)
	currRole := roles[0]

	r.GET("/api/role/:id", role.GetRoleDetail(svc))

	req, _ := http.NewRequest("GET", "/api/role/"+currRole.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := &responses.RoleResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, currRole.ID, resp.ID)
}

func TestCreateRoleHandler(t *testing.T) {
	r := SetUpRouter()
	svc := ROLESVC

	permissions, _, _ := PERMISSIONREPO.SelectPermission(context.Background(), nil)

	r.POST("/api/role", role.PostRole(svc))

	roleName := "testrole"
	payload := role.PostRolePayload{
		Name:        roleName,
		Description: "testing role",
	}
	for _, permission := range permissions {
		payload.Permissions = append(payload.Permissions, dto.IdParam{Id: permission.ID})
	}
	payloadJson, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/role", bytes.NewBuffer(payloadJson))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := &responses.RoleResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, roleName, resp.Name)
	assert.Equal(t, len(permissions), len(resp.Permissions))
}

func TestUpdateRoleHandler(t *testing.T) {
	r := SetUpRouter()
	svc := ROLESVC

	roles, _, _ := ROLEREPO.SelectRole(context.Background(), &model.Role{Name: "testrole"})
	currRole := roles[0]

	permissions, _, _ := PERMISSIONREPO.SelectPermission(context.Background(), &model.Permission{Method: "BUILDING"})

	r.PUT("/api/role/:id", role.PutRole(svc))

	roleName := "testchangerole"
	payload := role.PutRolePayload{
		Name: roleName,
	}
	for _, permission := range permissions {
		payload.Permissions = append(payload.Permissions, dto.IdParam{Id: permission.ID})
	}
	payloadJson, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", "/api/role/"+currRole.ID, bytes.NewBuffer(payloadJson))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := &responses.RoleResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, roleName, resp.Name)
	assert.Equal(t, len(permissions), len(resp.Permissions))
}

func TestDeleteRoleHandler(t *testing.T) {
	r := SetUpRouter()
	svc := ROLESVC

	roles, _, _ := ROLEREPO.SelectRole(context.Background(), &model.Role{Name: "testchangerole"})
	currRole := roles[0]

	r.DELETE("/api/role/:id", role.DeleteRole(svc))

	req, _ := http.NewRequest("DELETE", "/api/role/"+currRole.ID, nil)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := ""
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)

	afterDelRoles, _, _ := ROLEREPO.SelectRole(context.Background(), &model.Role{Name: "testchangerole"})
	assert.Equal(t, 0, len(afterDelRoles))
}
