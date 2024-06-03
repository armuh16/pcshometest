package logic

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"pcstakehometest/config"
	"pcstakehometest/module/auth/dto"
	userDto "pcstakehometest/module/user/dto"
	userLogic "pcstakehometest/module/user/logic"
	"pcstakehometest/static"

	"pcstakehometest/package/jwt"
	"pcstakehometest/package/logger"
	"pcstakehometest/utilities"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// AuthLogic
type IAuthLogic interface {
	Login(context.Context, *dto.LoginRequest, *gorm.DB) (*dto.Response, error)
}

type AuthLogic struct {
	fx.In
	Logger    *logger.LogRus
	UserLogic userLogic.IUserLogic
}

// NewLogic :
func NewLogic(AuthLogic AuthLogic) IAuthLogic {
	return &AuthLogic
}

// Login
func (l *AuthLogic) Login(ctx context.Context, reqData *dto.LoginRequest, tx *gorm.DB) (*dto.Response, error) {

	// Validate request data
	if err := reqData.Validate(); err != nil {
		l.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusBadRequest)
	}

	// Check exist user by email
	userDetail, err := l.UserLogic.Find(ctx, &userDto.FindRequest{
		Username: reqData.Username,
	})
	if err != nil {
		l.Logger.Error(err)
		return nil, err
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(userDetail.Password), []byte(reqData.Password)); err != nil {
		l.Logger.Error(err)
		return nil, utilities.ErrorRequest(errors.New(static.InvalidAccessLogin), http.StatusBadRequest)
	}

	// Generate uuid for user jwt
	uuid, err := uuid.NewV4()
	if err != nil {
		l.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	// Generate access and refresh token
	token, err := jwt.RequestToken(ctx, jwt.ClaimData{
		UserID: userDetail.ID,
		UUID:   uuid.String(),
	}, config.Get().Auth.Secret, time.Now().Add(config.Get().Auth.ExpireAccessTokenDuration).Unix(), time.Now().Add(config.Get().Auth.ExpireRefreshTokenDuration).Unix())
	if err != nil {
		l.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	return &dto.Response{
		Token:        token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}
