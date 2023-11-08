package media

import (
	"net/http"

	"github.com/ericmarcelinotju/gram/domain/media"
	response "github.com/ericmarcelinotju/gram/utils/http"
	"github.com/gin-gonic/gin"
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
