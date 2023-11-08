package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.com/firelogik/helios/domain/module/auth"
	"gitlab.com/firelogik/helios/domain/module/user"
	response "gitlab.com/firelogik/helios/utils/http"
)

// NewRoutesFactory create and returns a factory to create routes for the acknowledgement
func NewRoutesFactory(group *gin.RouterGroup) func(service auth.Service, userSvc user.Service) {
	authRoutesFactory := func(service auth.Service, userSvc user.Service) {
		group.POST("login", Login(service))

		group.POST("logout", Logout(service))

		group.POST("change-password", func(c *gin.Context) {
			user, newPassword, err := BindChangePassword(c)
			if err != nil {
				response.ResponseError(c, err, http.StatusUnprocessableEntity)
				return
			}

			if err = userSvc.UpdateUserPassword(c, user, newPassword); err != nil {
				response.ResponseError(c, err, http.StatusInternalServerError)
				return
			}

			response.ResponseSuccess(c, nil)
		})

		// User allowed to reset his/her password without old password when supplied with forgot password token
		group.POST("reset-password", func(c *gin.Context) {
			newPassword, forgotToken, err := BindResetPassword(c)
			if err != nil {
				response.ResponseError(c, err, http.StatusUnprocessableEntity)
				return
			}

			if err = service.ResetPassword(c, newPassword, forgotToken); err != nil {
				response.ResponseError(c, err, http.StatusInternalServerError)
				return
			}

			response.ResponseSuccess(c, nil)
		})

		// User forgot his/her password, system will send email containing link to change his/her password
		group.POST("forgot-password", ForgotPassword(service))
	}
	return authRoutesFactory
}
