package handlers

import (
	"errors"
	"graph-interview/internal/api/handlers/dto"
	api_error "graph-interview/internal/api/handlers/errors"
	"graph-interview/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateUserReq  true  "User registration data"
// @Success      201   {object}  dto.Response{data=dto.CreateUserResp}
// @Failure      400   {object}  dto.Response
// @Router       /v1/auth/register [post]
func Register(userSrv *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := dto.CreateUserReq{}
		if err := c.ShouldBindJSON(&req); err != nil {
			dto.Err(c, err)
			return
		}
		resp, err := userSrv.CreateUser(c, req)
		if err != nil {
			dto.Err(c, err)
			return
		}
		dto.Created(c, "user created", resp)
	}
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginUserReq  true  "Login credentials"
// @Success      200   {object}  dto.Response{data=dto.JWTResp}
// @Failure      401   {object}  dto.Response
// @Router       /v1/auth/login [post]
func Login(authSrv *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := dto.LoginUserReq{}
		if err := c.ShouldBindJSON(&req); err != nil {
			dto.Err(c, err)
			return
		}
		resp, err := authSrv.LoginUser(c, req)
		if err != nil {
			if errors.Is(err, api_error.ErrInvalidCredentials) {
				dto.ErrUnauthorized(c, err)
				return
			}
			dto.ErrInternal(c, err)
			return
		}
		dto.OK(c, "login successful", resp)
	}
}

// Logout godoc
// @Summary      Logout user
// @Description  Revoke current access token
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.Response
// @Failure      401  {object}  dto.Response
// @Router       /v1/auth/logout [post]
func Logout(authSrv *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, _ := c.Cookie("access_token")
		if tokenStr == "" {
			h := c.GetHeader("Authorization")
			if len(h) > 7 {
				tokenStr = h[7:]
			}
		}
		if tokenStr == "" {
			dto.ErrUnauthorized(c, api_error.ErrUnauthorized)
			return
		}

		if err := authSrv.RevokeTokenByString(c, tokenStr); err != nil {
			dto.ErrInternal(c, err)
			return
		}
		authSrv.ClearAuthCookies(c)
		dto.OK(c, "logged out", nil)
	}
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Use refresh token to get new access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RefreshTokenReq  true  "Refresh token"
// @Success      200   {object}  dto.Response{data=dto.JWTResp}
// @Failure      401   {object}  dto.Response
// @Router       /v1/auth/refresh [post]
func RefreshToken(authSrv *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := dto.RefreshTokenReq{}
		if err := c.ShouldBindJSON(&req); err != nil {
			// Try from cookie
			refreshStr, cookieErr := c.Cookie("refresh_token")
			if cookieErr != nil || refreshStr == "" {
				dto.Err(c, err)
				return
			}
			req.RefreshToken = refreshStr
		}

		resp, err := authSrv.RefreshToken(c, req.RefreshToken)
		if err != nil {
			dto.ErrUnauthorized(c, err)
			return
		}
		dto.OK(c, "token refreshed", resp)
	}
}

// GetProfile godoc
// @Summary      Get user profile
// @Description  Get the authenticated user's profile
// @Tags         user
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.Response{data=dto.UserProfileResp}
// @Failure      401  {object}  dto.Response
// @Failure      404  {object}  dto.Response
// @Router       /v1/user/profile [get]
func GetProfile(userSrv *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := getUserID(c)
		if err != nil {
			dto.ErrUnauthorized(c, api_error.ErrUnauthorized)
			return
		}

		resp, err := userSrv.GetProfile(c, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, dto.Response{Error: err.Error()})
			return
		}
		dto.OK(c, "profile retrieved", resp)
	}
}

func getUserID(c *gin.Context) (uint, error) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		return 0, api_error.ErrUnauthorized
	}
	id, err := strconv.ParseUint(userIDStr.(string), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
