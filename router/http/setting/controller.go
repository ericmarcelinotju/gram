package setting

import (
	"net/http"

	"github.com/ericmarcelinotju/gram/domain/module/setting"
	dto "github.com/ericmarcelinotju/gram/router/http/dto/responses"
	response "github.com/ericmarcelinotju/gram/utils/http"
	"github.com/gin-gonic/gin"
)

// GetSetting godoc
// @Summary     Get list of settings
// @Description Get list of settings
// @Tags        Setting
// @Accept      json
// @Produce     json
// @Success     200    {object}   response.SetResponse{data=dto.ListSettingResponse}
// @Router      /setting  [get]
// @Security    Auth
func GetSetting(service setting.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		settings, err := service.ReadSetting(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		result := dto.ListSettingResponse{
			Settings: make([]dto.SettingResponse, len(settings)),
		}

		for i, setting := range settings {
			result.Settings[i] = *setting.ToResponseModel()
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
// @Param       setting   body       SettingPayload   true   "Setting Data"
// @Success     200       {object}   response.SetResponse{data=dto.SettingResponse}
// @Router      /setting  [post]
// @Security    Auth
func SaveSetting(service setting.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		setting, err := BindSave(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		err = service.SaveSetting(c, setting)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, setting.ToResponseModel())
	}
}
