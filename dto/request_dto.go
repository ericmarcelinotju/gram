package dto

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IdDto struct {
	Id string `json:"id" form:"id" uri:"id" binding:"required,uuid"`
}

type PaginationDto struct {
	Limit *int `json:"limit" form:"limit" uri:"limit" binding:"omitempty,min=1"`
	Page  *int `json:"page" form:"page" uri:"page" binding:"omitempty,min=1"`
}

func (p PaginationDto) Apply(db *gorm.DB) *gorm.DB {
	if p.Limit != nil && p.Page != nil {

		offset := *p.Limit * (*p.Page - 1)
		if *p.Limit > 0 && offset >= 0 {
			return db.
				Limit(*p.Limit).
				Offset(offset)
		}
	}
	return db
}

type SortDto struct {
	Sort *string `json:"sort" form:"sort" uri:"sort"`
}

func (s SortDto) Apply(db *gorm.DB) *gorm.DB {
	sorts := strings.Split(*s.Sort, ":")
	if len(sorts) == 2 {
		column := sorts[0]
		order := sorts[1]

		if len(column) > 0 && len(order) > 0 {
			order := fmt.Sprintf("%s %s", column, order)
			return db.Order(order)
		}
	}
	return db.Order("created_at DESC")
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
