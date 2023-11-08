package dto

import (
	"errors"
	"mime/multipart"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/firelogik/helios/domain/model"
)

type IdParam struct {
	Id string `json:"id" form:"id" uri:"id" binding:"required,uuid"`
}

type Pagination struct {
	TotalPerPage *int `json:"limit" form:"limit" uri:"limit" binding:"omitempty,min=1"`
	CurrentPage  *int `json:"page" form:"page" uri:"page" binding:"omitempty,min=1"`
}

func ToPaginationModel(p Pagination) (*model.Pagination, error) {
	if p.TotalPerPage != nil && p.CurrentPage != nil {
		return &model.Pagination{
			Limit:  *p.TotalPerPage,
			Offset: *p.TotalPerPage * (*p.CurrentPage - 1),
		}, nil
	}
	return nil, errors.New("incomplete pagination data")
}

type Sort struct {
	Sort *string `json:"sort" form:"sort" uri:"sort"`
}

func ToSortModel(s Sort) (*model.Sort, error) {
	if s.Sort != nil {
		sorts := strings.Split(*s.Sort, ":")
		if len(sorts) == 2 {
			return &model.Sort{
				Column: sorts[0],
				Order:  sorts[1],
			}, nil
		}
	}

	return nil, errors.New("incomplete sort data")
}

func BindFile(c *gin.Context) (*multipart.File, error) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	return &file, nil
}
