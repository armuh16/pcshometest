package logic

import (
	"context"
	"fmt"
	"net/http"

	"pcstakehometest/model"
	"pcstakehometest/module/user/dto"
	"pcstakehometest/module/user/repository"
	"pcstakehometest/package/logger"
	"pcstakehometest/static"
	"pcstakehometest/utilities"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// UserLogic
type IUserLogic interface {
	Find(context.Context, *dto.FindRequest) (*model.Users, error)
}

type UserLogic struct {
	fx.In
	Logger   *logger.LogRus
	UserRepo repository.IUserRepository
}

// NewLogic :
func NewLogic(userLogic UserLogic) IUserLogic {
	return &userLogic
}

// FindByID
func (l *UserLogic) Find(ctx context.Context, reqData *dto.FindRequest) (*model.Users, error) {
	product, err := l.UserRepo.Find(ctx, &model.Users{
		ID:       reqData.ID,
		Username: reqData.Username,
	})
	if err != nil {
		l.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, utilities.ErrorRequest(fmt.Errorf(static.DataNotFound, "user"), http.StatusNotFound)
		}
		return nil, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}
	return product, nil
}
