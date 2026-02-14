package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Msg     string `json:"msg,omitempty"`
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func ErrStatus(c *gin.Context, status int, err error) {
	c.JSON(status, Response{Success: false, Error: err.Error()})
}

func Err(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Response{Success: false, Error: err.Error()})
}
