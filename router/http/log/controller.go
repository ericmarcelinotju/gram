package log

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/firelogik/helios/domain/module/log"
	dto "gitlab.com/firelogik/helios/router/http/dto/responses"
	response "gitlab.com/firelogik/helios/utils/http"
)

// GetLog godoc
// @Summary     Get list of logs
// @Description Get list of logs
// @Tags        Log
// @Accept      json
// @Produce     json
// @Param       item   query      LogFilter   true   "Paging, Search & Filter"
// @Success     200    {object}   response.SetResponse{data=dto.ListLogResponse}
// @Router      /log  [get]
// @Security    Auth
func GetLog(service log.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := BindGet(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}

		logs, err := service.ReadLog(c, payload)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		result := dto.ListLogResponse{
			Logs: make([]dto.LogResponse, len(logs)),
		}

		for i, log := range logs {
			result.Logs[i] = *log.ToResponseModel()
		}

		response.ResponseSuccess(c, result)
	}
}

// GetLogDetail godoc
// @Summary     Get log's detail
// @Description  Get log's detail
// @Tags        Log
// @Accept      json
// @Produce     json
// @Param       id    path       string   true   "Log ID"
// @Success     200   {object}   response.SetResponse{data=dto.LogResponse}
// @Router      /log/{id}  [get]
// @Security    Auth
func GetLogDetail(service log.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := BindID(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		log, err := service.ReadLogByID(c, id)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, log.ToResponseModel())
	}
}

// PostLog godoc
// @Summary     Post new log
// @Description Create new log
// @Tags        Log
// @Accept      json
// @Produce     json
// @Param       log   body       PostLogPayload   true   "Log Data"
// @Success     200          {object}   response.SetResponse{data=dto.LogResponse}
// @Router      /log  [post]
// @Security    Auth
func PostLog(service log.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		log, err := BindPost(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}
		err = service.CreateLog(c, log)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, log.ToResponseModel())
	}
}

// DeleteLog godoc
// @Summary     Delete log by id
// @Description Delete log by id
// @Tags        Log
// @Accept      json
// @Produce     json
// @Param       id    path       string   true   "Log ID"
// @Success     200   {object}   response.SetResponse
// @Router      /log/{id} [delete]
// @Security    Auth
func DeleteLog(service log.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := BindID(c)
		if err != nil {
			response.ResponseError(c, err, http.StatusUnprocessableEntity)
			return
		}

		err = service.DeleteLogByID(c, id)
		if err != nil {
			response.ResponseError(c, err, http.StatusInternalServerError)
			return
		}

		response.ResponseSuccess(c, nil)
	}
}
