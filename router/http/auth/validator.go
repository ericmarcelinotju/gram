package auth

import (
	"github.com/ericmarcelinotju/gram/domain/model"
	"github.com/gin-gonic/gin"
)

type LoginValidator struct {
	Username     string `json:"username" form:"username" binding:"required"`
	Password     string `json:"password" form:"password" binding:"required"`
	IsRememberMe bool   `json:"remember_me" form:"remember_me"`
}

func BindLogin(c *gin.Context) (*model.User, bool, error) {
	var payload LoginValidator
	if err := c.ShouldBind(&payload); err != nil {
		return nil, false, err
	}

	user := &model.User{
		Username: payload.Username,
		Password: payload.Password,
	}

	return user, payload.IsRememberMe, nil
}

type ChangePasswordValidator struct {
	UserID          string `form:"id" json:"id"`
	OldPassword     string `form:"old_password" json:"old_password"`
	NewPassword     string `form:"new_password" json:"new_password"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"eqfield=NewPassword"`
}

func BindChangePassword(c *gin.Context) (*model.User, string, error) {
	var payload ChangePasswordValidator
	if err := c.ShouldBind(&payload); err != nil {
		return nil, "", err
	}

	user := &model.User{
		ID:       payload.UserID,
		Password: payload.OldPassword,
	}

	return user, payload.NewPassword, nil
}

type ResetPasswordValidator struct {
	NewPassword     string `form:"new_password" json:"new_password"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"eqfield=NewPassword"`

	ForgotToken string `form:"forgot_token" json:"forgot_token"`
}

func BindResetPassword(c *gin.Context) (string, string, error) {
	var payload ResetPasswordValidator
	if err := c.ShouldBind(&payload); err != nil {
		return "", "", err
	}

	return payload.NewPassword, payload.ForgotToken, nil
}

type ForgotPasswordValidator struct {
	Username string `form:"username" json:"username" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,email"`
}

func BindForgotPassword(c *gin.Context) (*model.User, error) {
	var payload ForgotPasswordValidator
	if err := c.ShouldBind(&payload); err != nil {
		return nil, err
	}

	user := &model.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	return user, nil
}
