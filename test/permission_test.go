package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	"gitlab.com/firelogik/helios/domain/model"
	"gitlab.com/firelogik/helios/router/http/dto"
	responses "gitlab.com/firelogik/helios/router/http/dto/responses"
	"gitlab.com/firelogik/helios/router/http/permission"

	response "gitlab.com/firelogik/helios/utils/http"
)

func TestListPermissionHandler(t *testing.T) {
	r := SetUpRouter()
	svc := PERMISSIONSVC

	r.GET("/api/permission", permission.GetPermission(svc))
	req, _ := http.NewRequest("GET", "/api/permission", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := &responses.ListPermissionResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, 0, len(resp.Permissions))
}

func TestListPaginationPermissionHandler(t *testing.T) {
	r := SetUpRouter()
	svc := PERMISSIONSVC

	r.GET("/api/permission", permission.GetPermission(svc))

	currentPage := 1
	totalPerPage := 10

	request, _ := json.Marshal(permission.PermissionFilter{
		Pagination: dto.Pagination{
			CurrentPage:  &currentPage,
			TotalPerPage: &totalPerPage,
		},
	})
	req, _ := http.NewRequest("GET", "/api/permission", bytes.NewBuffer(request))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := &responses.ListPermissionResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, totalPerPage, len(resp.Permissions))
}

func TestDetailPermissionHandler(t *testing.T) {
	r := SetUpRouter()
	svc := PERMISSIONSVC

	permissions, _, _ := PERMISSIONREPO.SelectPermission(context.Background(), nil)
	currPermission := permissions[0]

	r.GET("/api/permission/:id", permission.GetPermissionDetail(svc))

	req, _ := http.NewRequest("GET", "/api/permission/"+currPermission.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := &responses.PermissionResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, currPermission.ID, resp.ID)
}

func TestCreatePermissionHandler(t *testing.T) {
	r := SetUpRouter()
	svc := PERMISSIONSVC

	r.POST("/api/permission", permission.PostPermission(svc))

	permissionModule := "testpermission"
	payload := permission.PostPermissionPayload{
		Module: permissionModule,
		Method: "POST",
	}
	payloadJson, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/permission", bytes.NewBuffer(payloadJson))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := &responses.PermissionResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, permissionModule, resp.Module)
}

func TestUpdatePermissionHandler(t *testing.T) {
	r := SetUpRouter()
	svc := PERMISSIONSVC

	permissions, _, _ := PERMISSIONREPO.SelectPermission(context.Background(), &model.Permission{Module: "testpermission"})
	currPermission := permissions[0]

	r.PUT("/api/permission/:id", permission.PutPermission(svc))

	permissionModule := "testchangepermission"
	payload := permission.PutPermissionPayload{
		Module: permissionModule,
		Method: "POST",
	}
	payloadJson, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", "/api/permission/"+currPermission.ID, bytes.NewBuffer(payloadJson))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := &responses.PermissionResponse{}
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, permissionModule, resp.Module)
}

func TestDeletePermissionHandler(t *testing.T) {
	r := SetUpRouter()
	svc := PERMISSIONSVC

	permissions, _, _ := PERMISSIONREPO.SelectPermission(context.Background(), &model.Permission{Module: "testchangepermission"})
	currPermission := permissions[0]

	r.DELETE("/api/permission/:id", permission.DeletePermission(svc))

	req, _ := http.NewRequest("DELETE", "/api/permission/"+currPermission.ID, nil)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := ""
	setResp := response.SetResponse{Data: resp}
	json.Unmarshal(w.Body.Bytes(), &setResp)

	assert.Equal(t, http.StatusOK, w.Code)

	afterDelPermissions, _, _ := PERMISSIONREPO.SelectPermission(context.Background(), &model.Permission{Module: "testchangepermission"})
	assert.Equal(t, 0, len(afterDelPermissions))
}
