package handlers

import (
	"graph-interview/internal/api/handlers/dto"
	"graph-interview/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 登录
// @Description 登录
// @Produce json
// @Param body body controllers.LoginParams true "body参数"
// @Success 200 {string} string "ok" "返回用户信息"
// @Failure 400 {string} string "err_code：10002 参数错误； err_code：10003 校验错误"
// @Failure 401 {string} string "err_code：10001 登录失败"
// @Failure 500 {string} string "err_code：20001 服务错误；err_code：20002 接口错误；err_code：20003 无数据错误；err_code：20004 数据库异常；err_code：20005 缓存异常"
// @Router /user/person/login [post
func CreateUser(userSrv *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := dto.CreateUserReq{}
		err := c.BindJSON(req)
		if err != nil {
			dto.Err(c, err)
			return
		}
		resp, err := userSrv.CreateUser(c, req)
		if err != nil {
			dto.Err(c, err)
			return
		}
		c.JSON(http.StatusCreated, &dto.Response{
			Msg:  "user created",
			Data: resp,
		})
	}
}

func Logout(authSrv *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := dto.CreateUserReq{}
		err := c.BindJSON(req)
		if err != nil {
			dto.Err(c, err)
			return
		}
		resp, err := authSrv.RevokeToken(c, req)
		if err != nil {
			dto.Err(c, err)
			return
		}
		c.JSON(http.StatusCreated, &dto.Response{
			Msg:  "user created",
			Data: resp,
		})
	}
}
