package http

import (
	"errors"
	"time"

	"github.com/ericmarcelinotju/gram/domain/model"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) (*model.User, error) {
	userCtx, ok := c.Get("auth-user")
	if !ok {
		return nil, errors.New("no user found in context")
	}
	user, ok := userCtx.(*model.User)
	if !ok || user == nil {
		return nil, errors.New("user context format invalid")
	}
	return user, nil
}

func authTokenHeaderLookup(c *gin.Context) *string {
	authHeader := c.GetHeader("authorization")
	if len(authHeader) <= 0 {
		return nil
	}
	return &authHeader
}

func authTokenCookieLookup(c *gin.Context) *string {
	authCookie, err := c.Request.Cookie("auth")
	if err != nil {
		return nil
	} else if authCookie.Expires.After(time.Now()) {
		return nil
	}
	return &authCookie.Value
}

func GetAuthToken(c *gin.Context) (string, error) {
	var token *string
	authHeader := authTokenHeaderLookup(c)
	token = authHeader

	// Cookie lookup
	if token == nil {
		authCookie := authTokenCookieLookup(c)
		token = authCookie
	}

	if token == nil {
		return "", errors.New("token not found")
	}
	return *token, nil
}
