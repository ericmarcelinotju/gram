package user

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"gitlab.com/firelogik/helios/domain/model"
	"gitlab.com/firelogik/helios/router/http/dto"
)

type PostUserPayload struct {
	Username        string                `json:"username" form:"username" binding:"required,min=2"`
	Firstname       string                `json:"first_name" form:"first_name"`
	Lastname        string                `json:"last_name" form:"last_name"`
	Department      string                `json:"department" form:"department"`
	Title           string                `json:"title" form:"title"`
	Email           string                `json:"email" form:"email" binding:"required,email"`
	Password        string                `json:"password" form:"password" binding:"required"`
	ConfirmPassword string                `json:"confirm_password" form:"confirm_password" binding:"eqfield=Password"`
	Avatar          *multipart.FileHeader `form:"avatar" swaggerignore:"true"`
	RoleID          string                `json:"role_id" form:"role_id" binding:"required,uuid"`
}

type PutUserPayload struct {
	Username        string                `json:"username" form:"username" binding:"min=2"`
	Firstname       string                `json:"first_name" form:"first_name"`
	Lastname        string                `json:"last_name" form:"last_name"`
	Department      string                `json:"department" form:"department"`
	Title           string                `json:"title" form:"title"`
	Email           string                `json:"email" form:"email" binding:"email"`
	Password        string                `json:"password" form:"password"`
	ConfirmPassword string                `json:"confirm_password" form:"confirm_password" binding:"eqfield=Password"`
	Avatar          *multipart.FileHeader `form:"avatar" swaggerignore:"true"`
	RoleID          string                `json:"role_id" form:"role_id" binding:"required,uuid"`
}

type UserFilter struct {
	ID       *string `json:"id" form:"id" uri:"id" binding:"omitempty,uuid"`
	Username *string `json:"username" form:"username" uri:"username"`
	Email    *string `json:"email" form:"email" uri:"email"`
	RoleID   *string `json:"role_id" form:"role_id" uri:"role_id" binding:"omitempty,uuid"`

	dto.Pagination
	dto.Sort
}

func BindPost(c *gin.Context) (*model.User, error) {
	var json PostUserPayload
	if err := c.ShouldBind(&json); err != nil {
		return nil, err
	}

	payload := &model.User{
		Username:   json.Username,
		Firstname:  json.Firstname,
		Lastname:   json.Lastname,
		Department: json.Department,
		Title:      json.Title,
		Email:      json.Email,
		Password:   json.Password,
		RoleID:     json.RoleID,
	}
	if json.Avatar != nil {
		file, err := json.Avatar.Open()
		if err != nil {
			return nil, err
		}
		payload.AvatarFile = &file
	}

	return payload, nil
}

func BindPut(c *gin.Context) (*model.User, error) {
	var idParam dto.IdParam
	err := c.ShouldBindUri(&idParam)
	if err != nil {
		return nil, err
	}

	var json PutUserPayload
	if err := c.ShouldBind(&json); err != nil {
		return nil, err
	}
	payload := &model.User{
		ID:         idParam.Id,
		Username:   json.Username,
		Firstname:  json.Firstname,
		Lastname:   json.Lastname,
		Department: json.Department,
		Title:      json.Title,
		Email:      json.Email,
		Password:   json.Password,
		RoleID:     json.RoleID,
	}

	return payload, nil
}

func BindGet(c *gin.Context) (*model.User, error) {
	var json UserFilter
	if err := c.ShouldBind(&json); err != nil {
		return nil, err
	}

	user := &model.User{}

	if json.ID != nil {
		user.ID = *json.ID
	}
	if json.Username != nil {
		user.Username = *json.Username
	}
	if json.Email != nil {
		user.Email = *json.Email
	}
	if json.RoleID != nil {
		user.RoleID = *json.RoleID
	}
	pagination, _ := dto.ToPaginationModel(json.Pagination)
	if pagination != nil {
		user.Pagination = *pagination
	}
	sort, _ := dto.ToSortModel(json.Sort)
	if sort != nil {
		user.Sort = *sort
	}

	return user, nil
}

func BindID(c *gin.Context) (string, error) {
	var json UserFilter
	if err := c.ShouldBindUri(&json); err != nil {
		return "", err
	}

	return *json.ID, nil
}
