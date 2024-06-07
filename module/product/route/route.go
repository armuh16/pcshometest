package route

import (
	"errors"
	"net/http"

	"pcstakehometest/database/postgres"
	"pcstakehometest/enum"
	"pcstakehometest/module/product/dto"
	"pcstakehometest/module/product/logic"
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
	Logic     logic.IProductLogic
	EchoRoute *router.Router
	Logger    *logger.LogRus
	Db        *postgres.DB
}

func NewRoute(h Handler, m ...echo.MiddlewareFunc) Handler {
	h.Route(m...)
	return h
}

func (h *Handler) Route(m ...echo.MiddlewareFunc) {
	product := h.EchoRoute.Group("/v1/product", m...)
	product.POST("", h.Create, h.EchoRoute.Authentication)
	product.GET("", h.FindAll, h.EchoRoute.Authentication)
	product.GET("/list", h.FindAllForBuyer, h.EchoRoute.Authentication)
}

// FindAllForBuyer
func (h *Handler) FindAllForBuyer(c echo.Context) error {
	var reqData = new(dto.FindAllRequest)

	_, ok := c.Request().Context().Value(jwt.InternalClaimData{}).(jwt.InternalClaimData)
	if !ok {
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
		})
	}

	if err := echo.QueryParamsBinder(c).
		Int("seller", (&reqData.SellerID)).
		BindError(); err != nil {
		h.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.BadRequest), http.StatusBadRequest),
		})
	}

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

// FindAll
func (h *Handler) FindAll(c echo.Context) error {
	var reqData = new(dto.FindAllRequest)

	data, ok := c.Request().Context().Value(jwt.InternalClaimData{}).(jwt.InternalClaimData)
	if !ok || data.Role != enum.RoleTypeSeller {
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
		})
	}
	reqData.SellerID = data.UserID

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

// Create
func (h *Handler) Create(c echo.Context) error {
	var reqData = new(dto.CreateRequest)

	data, ok := c.Request().Context().Value(jwt.InternalClaimData{}).(jwt.InternalClaimData)
	if !ok {
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
		})
	}

	reqData.SellerID = data.UserID
	reqData.RoleID = data.Role

	if err := c.Bind(reqData); err != nil {
		h.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.BadRequest), http.StatusBadRequest),
		})
	}

	tx := h.Db.Gorm.Begin()
	if err := h.Logic.Create(c.Request().Context(), reqData, tx); err != nil {
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
