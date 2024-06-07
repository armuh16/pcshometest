package route

import (
	"errors"
	"net/http"

	"pcstakehometest/database/postgres"
	"pcstakehometest/module/auth/dto"
	"pcstakehometest/module/auth/logic"
	"pcstakehometest/package/logger"
	"pcstakehometest/router"
	"pcstakehometest/static"
	"pcstakehometest/utilities"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type Handler struct {
	fx.In
	Logic     logic.IAuthLogic
	EchoRoute *router.Router
	Logger    *logger.LogRus
	Db        *postgres.DB
}

func NewRoute(h Handler, m ...echo.MiddlewareFunc) Handler {
	h.Route(m...)
	return h
}

func (h *Handler) Route(m ...echo.MiddlewareFunc) {
	auth := h.EchoRoute.Group("/v1/auth", m...)
	auth.POST("/login", h.Login)
}

// Login
func (h *Handler) Login(c echo.Context) error {
	var reqData = new(dto.LoginRequest)

	if err := c.Bind(reqData); err != nil {
		h.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			Error: utilities.ErrorRequest(errors.New(static.BadRequest), http.StatusBadRequest),
		})
	}

	tx := h.Db.Gorm.Begin()
	resp, err := h.Logic.Login(c.Request().Context(), reqData, tx)
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
		Data:   resp,
	})
}
