package setting

import (
	"net/http"

	"github.com/ericmarcelinotju/gram/dto"
	"github.com/ericmarcelinotju/gram/utils/request"
	"github.com/ericmarcelinotju/gram/utils/response"
	"github.com/gin-gonic/gin"
)

// GetSetting godoc
// @Summary     Get list of settings
// @Description Get list of settings
// @Tags        Setting
// @Accept      json
// @Produce     json
// @Success     200    {object}   response.SetResponse{data=dto.ListSettingDto}
// @Router      /setting  [get]
// @Security    Auth
func Get(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		settings, err := service.Read(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		result := dto.ListDto[dto.SettingDto]{
			Data: settings,
		}

		response.ResponseSuccess(c, result)
	}
}

// SaveSetting godoc
// @Summary     Post new setting
// @Description Create new setting
// @Tags        Setting
// @Accept      json
// @Produce     json
// @Param       setting   body       dto.PostSettingDto   true   "Setting Data"
// @Success     200       {object}   response.SetResponse
// @Router      /setting  [post]
// @Security    Auth
func Save(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := request.Bind[dto.PostSettingDto](c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		err = service.Save(c, payload)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, nil)
	}
}
