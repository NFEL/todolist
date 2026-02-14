package services

import (
	"context"
	"errors"
	"fmt"
	"graph-interview/internal/api/handlers/dto"
	api_error "graph-interview/internal/api/handlers/errors"
	"graph-interview/internal/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo  repository.UserRepo
	JwtSecret []byte
	redis     *redis.Client
}

func NewAuthService(userRepo repository.UserRepo, redis *redis.Client, jwtSecret []byte) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		redis:     redis,
		JwtSecret: jwtSecret,
	}
}

func (s *AuthService) LoginUser(ctx context.Context, req dto.LoginUserReq) (*dto.JWTResp, error) {
	user, err := s.UserRepo.GetByField(ctx, "username", req.Username)
	if err != nil {
		return nil, api_error.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, api_error.ErrInvalidCredentials
	}

	tokens, err := s.IssueTokens(fmt.Sprintf("%d", user.ID))
	if err != nil {
		return nil, err
	}

	if err := s.Persist(ctx, tokens); err != nil {
		return nil, err
	}

	return &dto.JWTResp{
		Access:    tokens.Access,
		Refresh:   tokens.Refresh,
		ExpiresAt: time.Now().Add(tokens.ExpAcc),
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenStr string) (*dto.JWTResp, error) {
	claims, err := s.ParseToken(refreshTokenStr)
	if err != nil {
		return nil, api_error.ErrInvalidToken
	}

	// Check if refresh token exists in Redis
	if cmd := s.redis.Get(ctx, "refresh:"+claims.ID); cmd.Err() != nil {
		return nil, api_error.ErrTokenRevoked
	}

	// Delete old refresh token
	s.redis.Del(ctx, "refresh:"+claims.ID)

	// Issue new tokens
	tokens, err := s.IssueTokens(claims.Subject)
	if err != nil {
		return nil, err
	}

	if err := s.Persist(ctx, tokens); err != nil {
		return nil, err
	}

	return &dto.JWTResp{
		Access:    tokens.Access,
		Refresh:   tokens.Refresh,
		ExpiresAt: time.Now().Add(tokens.ExpAcc),
	}, nil
}

type Tokens struct {
	Access   string
	Refresh  string
	JTIAcc   string
	JTIRef   string
	ExpAcc   time.Duration
	ExpRef   time.Duration
	UserID   string
	Issuer   string
	Audience string
}

func (s *AuthService) IssueTokens(userID string) (*Tokens, error) {
	now := time.Now().UTC()
	t := &Tokens{
		UserID:   userID,
		JTIAcc:   uuid.NewString(),
		JTIRef:   uuid.NewString(),
		ExpAcc:   15 * time.Minute,
		ExpRef:   7 * 24 * time.Hour,
		Issuer:   "jwt-todo-app",
		Audience: "jwt-todo-client",
	}
	ExpRefFromNow := now.Add(7 * 24 * time.Hour)
	ExpAccFromNow := now.Add(15 * time.Minute)

	acc := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ID:        t.JTIAcc,
		Issuer:    t.Issuer,
		Audience:  jwt.ClaimStrings{t.Audience},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(ExpAccFromNow),
	})

	ref := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ID:        t.JTIRef,
		Issuer:    t.Issuer,
		Audience:  jwt.ClaimStrings{t.Audience},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(ExpRefFromNow),
	})

	var signErr error
	t.Access, signErr = acc.SignedString(s.JwtSecret)
	if signErr != nil {
		return nil, signErr
	}
	t.Refresh, signErr = ref.SignedString(s.JwtSecret)
	if signErr != nil {
		return nil, signErr
	}
	return t, nil
}

func (s *AuthService) Persist(ctx context.Context, t *Tokens) error {
	if cmd := s.redis.Set(ctx, "access:"+t.JTIAcc, t.UserID, t.ExpAcc); cmd.Err() != nil {
		return cmd.Err()
	}
	if cmd := s.redis.Set(ctx, "refresh:"+t.JTIRef, t.UserID, t.ExpRef); cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (s *AuthService) SetAuthCookies(c *gin.Context, t *Tokens) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", t.Access, int(t.ExpAcc.Seconds()), "/", "", true, true)
	c.SetCookie("refresh_token", t.Refresh, int(t.ExpRef.Seconds()), "/", "", true, true)
}

func (s *AuthService) ClearAuthCookies(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
}

func (s *AuthService) RevokeTokenByString(ctx context.Context, tokenStr string) error {
	claims, err := s.ParseToken(tokenStr)
	if err != nil {
		return err
	}
	s.redis.Del(ctx, "access:"+claims.ID)
	return nil
}

func (s *AuthService) RevokeToken(ctx context.Context, t *Tokens) error {
	if cmd := s.redis.Del(ctx, "access:"+t.JTIAcc); cmd.Err() != nil {
		return cmd.Err()
	}
	if cmd := s.redis.Del(ctx, "refresh:"+t.JTIRef); cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (s *AuthService) ParseToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	token, err := parser.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
