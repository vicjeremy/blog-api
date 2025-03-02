package api

// The data model

import (
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"time"
)

type Blog struct {
	ID       uint           `json:"id" gorm:"primaryKey"`
	Title    string         `json:"title"`
	Content  string         `json:"content"`
	Category string         `json:"category"`
	Tags     datatypes.JSON `json:"tags" gorm:"type:jsonb"`
	CreateAt time.Time      `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt time.Time      `json:"updateAt" gorm:"autoUpdateTime"`
}

type JsonResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func ResponseJSON(c *gin.Context, status int, message string, data any) {
	response := JsonResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	c.JSON(status, response)
}
