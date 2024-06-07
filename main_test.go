package main_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	authRoute "pcstakehometest/module/auth/route"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"pcstakehometest/config"
	"pcstakehometest/database/postgres"
	"pcstakehometest/enum"
	"pcstakehometest/module"
	productRoute "pcstakehometest/module/product/route"
	transactionRoute "pcstakehometest/module/transaction/route"
	"pcstakehometest/package/jwt"
	"pcstakehometest/package/logger"
	"pcstakehometest/router"
)

func TestMain(m *testing.M) {
	config.SetConfig()
	app := fx.New(
		fx.Provide(router.NewRouter),
		fx.Provide(postgres.NewPostgres),
		fx.Provide(logger.NewLogRus),
		module.BundleRepository,
		module.BundleLogic,
		module.BundleRoute,
		fx.Invoke(NewRouteTest),
	)
	app.Start(context.Background())
	defer app.Stop(context.Background())
	m.Run()
}

type RouteTest struct {
	fx.In
	AuthHandler        authRoute.Handler
	TransactionHandler transactionRoute.Handler
	ProductHandler     productRoute.Handler
}

var r RouteTest

func NewRouteTest(routeTest RouteTest) {
	r = routeTest
}

func TestTransaction(t *testing.T) {
	t.Run("FailedCreateProduct", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
			"Name":"Testes Product",
    		"Description":"Failed Tested",
    		"Price":999999
		}`))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwt.InternalClaimData{
			UserID: 1,
			Role:   enum.RoleTypeSeller,
		})
		c.SetRequest(c.Request().WithContext(ctx))

		// Assertions
		if assert.NoError(t, r.ProductHandler.Create(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("SuccessCreateProduct", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
			"Name":"Testes Product",
    		"Description":"Success Tested",
    		"Price":999999
		}`))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwt.InternalClaimData{
			UserID: 1,
			Role:   enum.RoleTypeSeller,
		})
		c.SetRequest(c.Request().WithContext(ctx))

		// Assertions
		if assert.NoError(t, r.ProductHandler.Create(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedSellerCreateOrder", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
			"SellerID":1,
			"Items":[1,2]
		}`))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwt.InternalClaimData{
			UserID: 1,
			Role:   enum.RoleTypeSeller,
		})
		c.SetRequest(c.Request().WithContext(ctx))

		// Assertions
		if assert.NoError(t, r.TransactionHandler.CreateOrder(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("SuccessBuyerCreateOrder", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
			"SellerID":1,
			"Items":[1,2]
		}`))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwt.InternalClaimData{
			UserID: 1,
			Role:   enum.RoleTypeBuyer,
		})
		c.SetRequest(c.Request().WithContext(ctx))

		// Assertions
		if assert.NoError(t, r.TransactionHandler.CreateOrder(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("SuccessAcceptOrder", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
			"TransactionID":1
		}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwt.InternalClaimData{
			UserID: 1,
			Role:   enum.RoleTypeSeller,
		})
		c.SetRequest(c.Request().WithContext(ctx))

		// Assertions
		if assert.NoError(t, r.TransactionHandler.AcceptOrder(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedAcceptOrder", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/accept", strings.NewReader(`{
			"TransactionID":9999
		}`)) // Assuming 9999 is an invalid TransactionID
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwt.InternalClaimData{
			UserID: 1,
			Role:   enum.RoleTypeSeller,
		})
		c.SetRequest(c.Request().WithContext(ctx))

		// Assertions
		if assert.NoError(t, r.TransactionHandler.AcceptOrder(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})

	t.Run("SuccessFindAll", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwt.InternalClaimData{
			UserID: 2,
			Role:   enum.RoleTypeBuyer,
		})
		c.SetRequest(c.Request().WithContext(ctx))

		// Assertions
		if assert.NoError(t, r.TransactionHandler.FindAll(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedFindAllUnauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		// No valid context set, simulating unauthorized access
		c.SetRequest(c.Request())

		// Assertions
		if assert.NoError(t, r.TransactionHandler.FindAll(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	})

	t.Run("SuccessHistory", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/history", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwt.InternalClaimData{
			UserID: 2,
			Role:   enum.RoleTypeBuyer,
		})
		c.SetRequest(c.Request().WithContext(ctx))

		// Assertions
		if assert.NoError(t, r.TransactionHandler.History(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedHistoryUnauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/history", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		// Simulating unauthorized access by setting invalid role
		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwt.InternalClaimData{
			UserID: 2,
			Role:   enum.RoleTypeSeller,
		})
		c.SetRequest(c.Request().WithContext(ctx))

		// Assertions
		if assert.NoError(t, r.TransactionHandler.History(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	})

	t.Run("SuccessLogin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(`{
			"Username":"Seller",
			"Password":"secret"
		}`))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.NoError(t, r.AuthHandler.Login(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedLogin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(`{
			"Username":"Seller",
			"Password":"wrongpassword"
		}`))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.NoError(t, r.AuthHandler.Login(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("SuccessFindAll", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwt.InternalClaimData{
			UserID: 1,
			Role:   enum.RoleTypeSeller,
		})
		c.SetRequest(c.Request().WithContext(ctx))

		// Assertions
		if assert.NoError(t, r.ProductHandler.FindAll(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedFindAllUnauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		// No valid context set, simulating unauthorized access
		c.SetRequest(c.Request())

		// Assertions
		if assert.NoError(t, r.ProductHandler.FindAll(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	})

	t.Run("SuccessFindAllForBuyer", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/list?seller=1", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwt.InternalClaimData{
			UserID: 2,
			Role:   enum.RoleTypeBuyer,
		})
		c.SetRequest(c.Request().WithContext(ctx))

		// Assertions
		if assert.NoError(t, r.ProductHandler.FindAllForBuyer(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedFindAllForBuyerUnauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/list?seller=1", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		// No valid context set, simulating unauthorized access
		c.SetRequest(c.Request())

		// Assertions
		if assert.NoError(t, r.ProductHandler.FindAllForBuyer(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	})
}
