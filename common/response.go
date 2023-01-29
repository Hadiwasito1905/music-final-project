package common

import "github.com/gin-gonic/gin"

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  []string    `json:"errors"`
	Data    interface{} `json:"data"`
}

func SendResponse(c *gin.Context, statusCode int, status bool, message string, errors []string, data interface{}) {
	c.JSON(statusCode, Response{status, message, errors, data})
}
