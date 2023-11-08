package auth

import (
	"net/http"
	"strings"

	"github.com/ericmarcelinotju/gram/domain/module/auth"
	dto "github.com/ericmarcelinotju/gram/router/http/dto/responses"
	response "github.com/ericmarcelinotju/gram/utils/http"
	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary     Login
// @Description Login using email and password to generate token for auth
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       credential   body      LoginValidator   true   "Login Credential"
// @Success     200          {object}  response.SetResponse{data=dto.LoginResponse}
// @Router      /auth/login  [post]
func Login(service auth.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		user, isRememberMe, err := BindLogin(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}
		token, err := service.Login(c, user, isRememberMe)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, dto.LoginResponse{
			Token: token,
			User:  *user.ToResponseModel(),
		})
	}
}

// Logout godoc
// @Summary     Logout
// @Description Logout the current user determined by it's auth token
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Success     200    {object}   response.SetResponse{data=string}
// @Router      /auth/logout  [post]
// @Security    Auth
func Logout(service auth.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		token, err := response.GetAuthToken(c)
		if err != nil {
			response.ResponseAbort(c, err, http.StatusUnauthorized)
			return
		}
		err = service.Logout(c, token)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, nil)
	}
}

// ForgotPassword godoc
// @Summary     ForgotPassword
// @Description ForgotPassword using email to generate token for reset password
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       credential   body      ForgotPasswordValidator   true   "Login Credential"
// @Success     200          {object}  response.SetResponse
// @Router      /auth/forgot-password  [post]
func ForgotPassword(service auth.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		user, err := BindForgotPassword(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}

		if err = service.ForgotPassword(c, user); err != nil {
			if strings.Contains(err.Error(), "NotFound") {
				response.ResponseError(c, err, http.StatusNotFound)
				return
			}
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, nil)
	}
}
