package dto

type LoginRespDto struct {
	Token string  `json:"token"`
	User  UserDto `json:"user"`
}

type LoginDto struct {
	Username     string `json:"username" form:"username" binding:"required"`
	Password     string `json:"password" form:"password" binding:"required"`
	IsRememberMe bool   `json:"remember_me" form:"remember_me"`
}

type ChangeUserPasswordDto struct {
	Id              string `json:"id" form:"id" uri:"id" binding:"required,uuid"`
	OldPassword     string `json:"old_password" form:"old_password"`
	NewPassword     string `json:"new_password" form:"new_password"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"eqfield=NewPassword"`
}

type ResetUserPasswordDto struct {
	NewPassword     string `form:"new_password" json:"new_password"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"eqfield=NewPassword"`
	ForgotToken     string `form:"forgot_token" json:"forgot_token"`
}

type ForgotUserPasswordDto struct {
	Username string `form:"username" json:"username" binding:"required"`
}
