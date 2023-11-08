package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gitlab.com/firelogik/helios/domain/model"
	"gitlab.com/firelogik/helios/domain/module/auth"
	response "gitlab.com/firelogik/helios/utils/http"
)

type AuthMiddleware struct {
	Authenticate gin.HandlerFunc
	Authorize    gin.HandlerFunc
}

// NewRoutesFactory create and returns a factory to create routes for the acknowledgement
func NewAuthMiddleware(authSvc auth.Service) AuthMiddleware {
	getPermission := func(permissions []model.Permission, url, method string) *model.Permission {
		for _, perm := range permissions {
			if strings.Contains(url, strings.ToLower(perm.Module)) && perm.Method == method {
				return &perm
			}
		}
		return nil
	}

	getSpecialPermission := func(permissions []model.Permission, module, method string) *model.Permission {
		for _, perm := range permissions {
			if strings.ToLower(module) == strings.ToLower(perm.Module) && perm.Method == method {
				return &perm
			}
		}
		return nil
	}

	authMiddleware := AuthMiddleware{
		// Authenticate middleware
		Authenticate: func(c *gin.Context) {
			token, err := response.GetAuthToken(c)
			if err != nil {
				response.ResponseAbort(c, err, http.StatusUnauthorized)
				return
			}

			user, err := authSvc.ReadUserByToken(c, token)
			if err != nil {
				response.ResponseAbort(c, err, http.StatusUnauthorized)
				return
			}

			c.Set("auth-user", user)
			c.Set("auth-token", token)

			c.Next()
		},
		// Authorize middleware
		Authorize: func(c *gin.Context) {
			userCtx, ok := c.Get("auth-user")
			if !ok {
				response.ResponseAbort(c, errors.New("no user found in context"), http.StatusUnauthorized)
				return
			}
			user, ok := userCtx.(*model.User)
			if !ok || user == nil {
				response.ResponseAbort(c, errors.New("user context format invalid"), http.StatusUnauthorized)
				return
			}

			var permissions []model.Permission = user.Role.Permissions
			var urlString string = c.Request.URL.String()
			var requestMethod string = c.Request.Method

			// Special permissions
			if strings.Contains(urlString, "branch/sync") {
				if permission := getSpecialPermission(permissions, "SYNCHRONIZE_BRANCH", requestMethod); permission == nil {
					c.Set("permission", permission)
				} else {
					response.ResponseAbort(c, errors.New("user have no SYNCHRONIZE_BRANCH permission"), http.StatusForbidden)
					return
				}
			} else {
				if permission := getPermission(permissions, urlString, requestMethod); permission != nil {
					c.Set("permission", permission)
				} else {
					response.ResponseAbort(c, errors.New("user have no permission"), http.StatusForbidden)
					return
				}
			}

			c.Next()
		},
	}
	return authMiddleware
}
