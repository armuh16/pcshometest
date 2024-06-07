package route

import (
	"errors"
	"net/http"
	"pcstakehometest/enum"

	"pcstakehometest/database/postgres"
	"pcstakehometest/module/transaction/dto"
	"pcstakehometest/module/transaction/logic"
	"pcstakehometest/package/jwt"
	"pcstakehometest/package/logger"
	"pcstakehometest/router"
	"pcstakehometest/static"
	"pcstakehometest/utilities"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type Handler struct {
	fx.In
	Logic     logic.ITransactionLogic
	EchoRoute *router.Router
	Logger    *logger.LogRus
	Db        *postgres.DB
}

func NewRoute(h Handler, m ...echo.MiddlewareFunc) Handler {
	h.Route(m...)
	return h
}

func (h *Handler) Route(m ...echo.MiddlewareFunc) {
	transaction := h.EchoRoute.Group("/v1/transaction", m...)
	transaction.GET("", h.FindAll, h.EchoRoute.Authentication)
	transaction.POST("", h.CreateOrder, h.EchoRoute.Authentication)
	transaction.POST("", h.AcceptOrder, h.EchoRoute.Authentication)
	transaction.GET("/history", h.History, h.EchoRoute.Authentication)
}

// FindAll
func (h *Handler) FindAll(c echo.Context) error {
	var reqData = new(dto.FindAllRequest)

	data, ok := c.Request().Context().Value(jwt.InternalClaimData{}).(jwt.InternalClaimData)
	if !ok {
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
		})
	}

	reqData.UserID = data.UserID
	reqData.RoleID = data.Role

	resp, err := h.Logic.FindAll(c.Request().Context(), reqData)
	if err != nil {
		h.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: err,
		})
	}

	return utilities.Response(c, &utilities.ResponseRequest{
		Code:   http.StatusOK,
		Status: static.Success,
		Data:   resp,
	})
}

// CreateOrder
func (h *Handler) CreateOrder(c echo.Context) error {
	var reqData = new(dto.CreateOrderRequest)

	if err := c.Bind(reqData); err != nil {
		h.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.BadRequest), http.StatusBadRequest),
		})
	}

	data, ok := c.Request().Context().Value(jwt.InternalClaimData{}).(jwt.InternalClaimData)
	if !ok {
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
		})
	}

	reqData.BuyerID = data.UserID
	reqData.RoleID = data.Role

	tx := h.Db.Gorm.Begin()
	coupons, err := h.Logic.CreateOrder(c.Request().Context(), reqData, tx)
	if err != nil {
		h.Logger.Error(err)
		defer func() {
			tx.Rollback()
		}()
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: err,
		})
	}
	tx.Commit()

	return utilities.Response(c, &utilities.ResponseRequest{
		Code:   http.StatusOK,
		Status: static.Success,
		Data: map[string]interface{}{
			"coupons": coupons,
		},
	})
}

// AcceptOrder
func (h *Handler) AcceptOrder(c echo.Context) error {
	var reqData = new(dto.AcceptOrderRequest)

	if err := c.Bind(reqData); err != nil {
		h.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.BadRequest), http.StatusBadRequest),
		})
	}

	data, ok := c.Request().Context().Value(jwt.InternalClaimData{}).(jwt.InternalClaimData)
	if !ok {
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
		})
	}

	reqData.SellerID = data.UserID
	reqData.RoleID = data.Role

	tx := h.Db.Gorm.Begin()
	if err := h.Logic.AcceptOrder(c.Request().Context(), reqData, tx); err != nil {
		h.Logger.Error(err)
		defer func() {
			tx.Rollback()
		}()
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: err,
		})
	}
	tx.Commit()

	return utilities.Response(c, &utilities.ResponseRequest{
		Code:   http.StatusOK,
		Status: static.Success,
	})
}

// GetCouponsAndHistory
func (h *Handler) History(c echo.Context) error {
	var reqData = new(dto.FindHistory)

	data, ok := c.Request().Context().Value(jwt.InternalClaimData{}).(jwt.InternalClaimData)
	if !ok || data.Role != enum.RoleTypeBuyer {
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
		})
	}

	reqData.UserID = data.UserID
	reqData.RoleID = data.Role

	transactions, coupons, err := h.Logic.FindHistory(c.Request().Context(), reqData)
	if err != nil {
		h.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: err,
		})
	}

	return utilities.Response(c, &utilities.ResponseRequest{
		Code:   http.StatusOK,
		Status: static.Success,
		Data: map[string]interface{}{
			"transactions": transactions,
			"coupons":      coupons,
		},
	})
}
