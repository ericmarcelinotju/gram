package setting

import (
	"net/http"

	"github.com/ericmarcelinotju/gram/dto"
	httpUtil "github.com/ericmarcelinotju/gram/utils/http"
	"github.com/gin-gonic/gin"
)

// GetSetting godoc
// @Summary     Get list of settings
// @Description Get list of settings
// @Tags        Setting
// @Accept      json
// @Produce     json
// @Success     200    {object}   httpUtil.SetResponse{data=dto.ListSettingResponse}
// @Router      /setting  [get]
// @Security    Auth
func Get(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		settings, err := service.Read(c)
		if err != nil {
			httpUtil.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		result := dto.ListSettingDto{
			Settings: settings,
		}

		httpUtil.ResponseSuccess(c, result)
	}
}

// SaveSetting godoc
// @Summary     Post new setting
// @Description Create new setting
// @Tags        Setting
// @Accept      json
// @Produce     json
// @Param       setting   body       SettingPayload   true   "Setting Data"
// @Success     200       {object}   httpUtil.SetResponse{data=dto.SettingResponse}
// @Router      /setting  [post]
// @Security    Auth
func Save(service Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := httpUtil.Bind[dto.PostSettingDto](c)
		if err != nil {
			httpUtil.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		err = service.Save(c, payload)
		if err != nil {
			httpUtil.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		httpUtil.ResponseSuccess(c, nil)
	}
}
