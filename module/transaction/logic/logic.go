package logic

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"pcstakehometest/enum"
	"pcstakehometest/model"
	productDto "pcstakehometest/module/product/dto"
	productLogic "pcstakehometest/module/product/logic"
	"pcstakehometest/module/transaction/dto"
	"pcstakehometest/module/transaction/repository"
	userDto "pcstakehometest/module/user/dto"
	userLogic "pcstakehometest/module/user/logic"
	"pcstakehometest/package/logger"
	"pcstakehometest/static"
	"pcstakehometest/utilities"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// TransactionLogic
type ITransactionLogic interface {
	CreateOrder(context.Context, *dto.CreateOrderRequest, *gorm.DB) (int, error)
	FindAll(context.Context, *dto.FindAllRequest) ([]*model.Transactions, error)
	AcceptOrder(context.Context, *dto.AcceptOrderRequest, *gorm.DB) error
	FindHistory(context.Context, *dto.FindHistory) ([]*dto.TransactionHistoryResponse, int, error)
}

type TransactionLogic struct {
	fx.In
	Logger          *logger.LogRus
	ProductLogic    productLogic.IProductLogic
	UserLogic       userLogic.IUserLogic
	TransactionRepo repository.ITransactionRepository
}

// NewLogic :
func NewLogic(transactionLogic TransactionLogic) ITransactionLogic {
	return &transactionLogic
}

// Create
func (l *TransactionLogic) CreateOrder(ctx context.Context, reqData *dto.CreateOrderRequest, tx *gorm.DB) (int, error) {
	// Validate request data
	if err := reqData.Validate(); err != nil {
		l.Logger.Error(err)
		//return utilities.ErrorRequest(err, http.StatusBadRequest)
		return 0, utilities.ErrorRequest(err, http.StatusBadRequest)
	}

	var (
		grandTotal   float64
		snapshotItem model.ItemsTransaction
		coupons      int
	)

	// Find detail seller
	sellerDetail, err := l.UserLogic.Find(ctx, &userDto.FindRequest{
		ID: reqData.SellerID,
	})
	if err != nil {
		l.Logger.Error(err)
		return 0, err
	}

	// Find detail buyer
	buyerDetail, err := l.UserLogic.Find(ctx, &userDto.FindRequest{
		ID: reqData.BuyerID,
	})
	if err != nil {
		l.Logger.Error(err)
		return 0, err
	}

	// Validate product
	for _, productID := range reqData.Items {
		productDetail, err := l.ProductLogic.Find(ctx, &productDto.FindRequest{
			ID:       productID,
			SellerID: sellerDetail.ID,
		})
		if err != nil {
			l.Logger.Error(err)
			return 0, err
		}

		snapshotItem = append(snapshotItem, model.ProductTransaction{
			ID:          productDetail.ID,
			Name:        productDetail.Name,
			Description: productDetail.Description,
			Price:       productDetail.Price,
		})

		grandTotal += productDetail.Price

		if productDetail.Price > 49999 && productDetail.Price < 99999 {
			coupons++
		}

	}

	coupons += int(grandTotal) / 100000

	if _, err := l.TransactionRepo.Create(ctx, &model.Transactions{
		BuyerID:    buyerDetail.ID,
		SellerID:   sellerDetail.ID,
		GrandTotal: grandTotal,
		Status:     enum.TransactionStatusTypePending,
		Items:      snapshotItem,
	}, tx); err != nil {
		return 0, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	return coupons, nil
}

// GetCouponsAndHistory
func (l *TransactionLogic) FindHistory(ctx context.Context, reqData *dto.FindHistory) ([]*dto.TransactionHistoryResponse, int, error) {
	// Validate request data
	if err := reqData.Validate(); err != nil {
		l.Logger.Error(err)
		return nil, 0, utilities.ErrorRequest(err, http.StatusBadRequest)
	}

	var (
		whereData = model.Transactions{
			Status: enum.TransactionStatusTypeAccept,
		}
		coupons int
	)

	if reqData.RoleID == enum.RoleTypeBuyer {
		whereData.BuyerID = reqData.UserID
	} else {
		whereData.SellerID = reqData.UserID
	}

	transactions, err := l.TransactionRepo.FindHistory(ctx, &whereData)
	if err != nil {
		l.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, 0, utilities.ErrorRequest(fmt.Errorf(static.DataNotFound, "transaksi"), http.StatusNotFound)
		}
		return nil, 0, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	response := make([]*dto.TransactionHistoryResponse, len(transactions))
	for i, transaction := range transactions {
		coupons += transaction.Coupons
		transaction.StatusTransaction = transaction.Status.String()

		status := "Open"
		if time.Since(transaction.CreatedAt) > 3*time.Hour {
			status = "Closed"
		}

		response[i] = &dto.TransactionHistoryResponse{
			Transaction: *transaction,
			Status:      status,
		}
	}

	return response, coupons, nil
}

// FindAll
func (l *TransactionLogic) FindAll(ctx context.Context, reqData *dto.FindAllRequest) ([]*model.Transactions, error) {
	// Validate request data
	if err := reqData.Validate(); err != nil {
		l.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusBadRequest)
	}

	var whereData = model.Transactions{}

	if reqData.RoleID == enum.RoleTypeBuyer {
		whereData = model.Transactions{
			BuyerID: reqData.UserID,
		}
	} else {
		whereData = model.Transactions{
			SellerID: reqData.UserID,
		}
	}

	transactions, err := l.TransactionRepo.FindAll(ctx, &whereData)
	if err != nil {
		l.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, utilities.ErrorRequest(fmt.Errorf(static.DataNotFound, "transaksi"), http.StatusNotFound)
		}
		return nil, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	for _, v := range transactions {
		v.StatusTransaction = v.Status.String()
	}

	return transactions, nil
}

// AcceptOrder
func (l *TransactionLogic) AcceptOrder(ctx context.Context, reqData *dto.AcceptOrderRequest, tx *gorm.DB) error {
	// Validate request data
	if err := reqData.Validate(); err != nil {
		l.Logger.Error(err)
		return utilities.ErrorRequest(err, http.StatusBadRequest)
	}

	transaction, err := l.TransactionRepo.Find(ctx, &model.Transactions{
		ID:       reqData.TransactionID,
		SellerID: reqData.SellerID,
	})
	if err != nil {
		l.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return utilities.ErrorRequest(fmt.Errorf(static.DataNotFound, "transaksi"), http.StatusNotFound)
		}
		return utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	var newCoupons int
	for _, item := range transaction.Items {
		if item.Price > 50000 && item.Price < 100000 {
			newCoupons++
		}
	}

	newCoupons += int(transaction.GrandTotal) / 100000

	if err := l.TransactionRepo.Update(ctx, &model.Transactions{
		ID:       reqData.TransactionID,
		SellerID: reqData.SellerID,
		Status:   enum.TransactionStatusTypeAccept,
		Coupons:  newCoupons,
	}, tx); err != nil {
		l.Logger.Error(err)
		return utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	return nil
}
