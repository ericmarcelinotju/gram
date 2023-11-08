package setting

import (
	"github.com/ericmarcelinotju/gram/domain/model"
	"github.com/gin-gonic/gin"
)

type SettingPayload struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func BindSave(c *gin.Context) (*model.Setting, error) {
	var json SettingPayload
	if err := c.ShouldBind(&json); err != nil {
		return nil, err
	}

	setting := &model.Setting{
		Name:  json.Name,
		Value: json.Value,
	}

	return setting, nil
}
