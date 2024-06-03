package repository

import (
	"context"

	"pcstakehometest/database/postgres"
	"pcstakehometest/model"
	"pcstakehometest/package/logger"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// SellerRepository
type ISellerRepository interface {
	Create(context.Context, *model.Products, *gorm.DB) (*int, error)
	FindAll(context.Context, *model.Products) ([]*model.Products, error)
	Find(context.Context, *model.Products) (*model.Products, error)
}

type SellerRepository struct {
	fx.In
	Logger   *logger.LogRus
	Database *postgres.DB
}

// NewRepository :
func NewRepository(sellerRepository SellerRepository) ISellerRepository {
	return &sellerRepository
}

// Create
func (l *SellerRepository) Create(ctx context.Context, reqData *model.Products, tx *gorm.DB) (*int, error) {
	if err := tx.WithContext(ctx).Create(&reqData).Error; err != nil {
		l.Logger.Error(err)
		return nil, err
	}
	return &reqData.ID, nil
}

// FindAll
func (l *SellerRepository) FindAll(ctx context.Context, reqData *model.Products) ([]*model.Products, error) {
	products := []*model.Products{}

	if err := l.Database.Gorm.WithContext(ctx).Model(&model.Products{}).
		Preload("Seller").
		Where(&model.Products{
			SellerID: reqData.SellerID,
		}).
		Order("id desc").
		Find(&products).
		Error; err != nil {
		l.Logger.Error(err)
		return nil, err
	}

	return products, nil
}

// Find
func (l *SellerRepository) Find(ctx context.Context, reqData *model.Products) (*model.Products, error) {
	product := new(model.Products)
	if err := l.Database.Gorm.WithContext(ctx).
		Where(&model.Products{
			ID:       reqData.ID,
			SellerID: reqData.SellerID,
		}).First(&product).Error; err != nil {
		l.Logger.Error(err)
		return nil, err
	}
	return product, nil
}
