package auth

import (
	"net/http"
	"strings"

	"github.com/ericmarcelinotju/gram/dto"
	"github.com/ericmarcelinotju/gram/module/user"
	"github.com/ericmarcelinotju/gram/utils/request"
	"github.com/ericmarcelinotju/gram/utils/response"
	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary     Login
// @Description Login using email and password to generate token for auth
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       credential   body      dto.LoginDto   true   "Login Credential"
// @Success     200          {object}  response.SetResponse
// @Router      /auth/login  [post]
func Login(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := request.Bind[dto.LoginDto](c)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}
		user, token, err := service.Login(c, payload)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, gin.H{
			"token": token,
			"user":  user,
		})
	}
}

// Logout godoc
// @Summary     Logout
// @Description Logout the current user determined by it's auth token
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Success     200    {object}   response.SetResponse
// @Router      /auth/logout  [post]
// @Security    Auth
func Logout(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		token, err := request.GetAuthToken(c)
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

// ChangePassword godoc
// @Summary     ChangePassword
// @Description ChangePassword user password
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       credential   body      dto.ChangeUserPasswordDto   true   "Login Credential"
// @Success     200          {object}  response.SetResponse
// @Router      /auth/change-password  [post]
func ChangePassword(userSvc user.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := request.Bind[dto.ChangeUserPasswordDto](c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}

		if err = userSvc.UpdatePassword(c, payload); err != nil {
			// TODO : process custom error
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

// ForgotPassword godoc
// @Summary     ForgotPassword
// @Description ForgotPassword using email to generate token for reset password
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       credential   body      dto.ForgotUserPasswordDto   true   "Login Credential"
// @Success     200          {object}  response.SetResponse
// @Router      /auth/forgot-password  [post]
func ForgotPassword(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := request.Bind[dto.ForgotUserPasswordDto](c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}

		if err = service.ForgotPassword(c, payload); err != nil {
			// TODO : process custom error
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

// ResetPassword godoc
// @Summary     ResetPassword
// @Description ResetPassword using forgot token
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       credential   body      dto.ResetUserPasswordDto   true   "Login Credential"
// @Success     200          {object}  response.SetResponse
// @Router      /auth/reset-password  [post]
func ResetPassword(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := request.Bind[dto.ResetUserPasswordDto](c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}

		if err = service.ResetPassword(c, payload); err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, nil)
	}
}
