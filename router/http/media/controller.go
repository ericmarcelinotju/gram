package media

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/firelogik/helios/domain/media"
	response "gitlab.com/firelogik/helios/utils/http"
)

// UploadFile godoc
// @Summary     Upload file
// @Description Upload file
// @Tags        Media
// @Accept      mpfd
// @Produce     json
// @Param       file   body       FilePayload   true   "File Data"
// @Success     200    {object}   response.SetResponse
// @Router      /media  [post]
// @Security    Auth
func UploadFile(service media.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		file, err := Bind(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}
		filename, err := service.Upload(c, file)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, filename)
	}
}
