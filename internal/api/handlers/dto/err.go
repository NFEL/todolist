package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Msg     string `json:"msg,omitempty"`
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
}

func OK(c *gin.Context, msg string, data any) {
	c.JSON(http.StatusOK, Response{Success: true, Msg: msg, Data: data})
}

func Created(c *gin.Context, msg string, data any) {
	c.JSON(http.StatusCreated, Response{Success: true, Msg: msg, Data: data})
}

func ErrStatus(c *gin.Context, status int, err error) {
	c.JSON(status, Response{Success: false, Error: err.Error()})
}

func Err(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Response{Success: false, Error: err.Error()})
}

func ErrUnauthorized(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, Response{Success: false, Error: err.Error()})
}

func ErrNotFound(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, Response{Success: false, Error: err.Error()})
}

func ErrInternal(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, Response{Success: false, Error: err.Error()})
}
