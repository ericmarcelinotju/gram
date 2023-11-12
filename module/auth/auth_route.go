package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/ericmarcelinotju/gram/module/user"
)

// NewRoutesFactory create and returns a factory to create routes for the acknowledgement
func NewRoutesFactory(router *gin.Engine) func(service Service, userSvc user.Service) {
	group := router.Group("/api/auth")

	authRoutesFactory := func(service Service, userSvc user.Service) {
		group.POST("login", Login(service))

		group.POST("logout", Logout(service))

		group.POST("change-password", ChangePassword(userSvc))

		// User allowed to reset his/her password without old password when supplied with forgot password token
		group.POST("reset-password", ResetPassword(service))

		// User forgot his/her password, system will send email containing link to change his/her password
		group.POST("forgot-password", ForgotPassword(service))
	}
	return authRoutesFactory
}
