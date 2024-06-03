package router

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	userRepo "pcstakehometest/module/user/repository"
	"pcstakehometest/package/logger"
	"pcstakehometest/static"
	"pcstakehometest/utilities"
)

type Router struct {
	*echo.Echo
	userRepo userRepo.IUserRepository
}

var RouteLog *logger.LogRus

func NewRouter(logger *logger.LogRus,
	userRepo userRepo.IUserRepository) *Router {

	e := echo.New()

	//Set middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			RouteLog = logger.Request()

			requestLog := RouteLog.WithFields(logrus.Fields{
				"path":       c.Request().URL.Path,
				"method":     c.Request().Method,
				"version":    c.Request().Header.Get("Version"),
				"queryParam": c.QueryParams(),
			})

			//If content-type not application/json will skip
			if c.Request().Header.Get("Content-Type") != echo.MIMEApplicationJSON {
				requestLog.Info()
				return next(c)
			}

			//Read reqData from body
			reqData, err := ioutil.ReadAll(c.Request().Body)
			if err != nil {
				RouteLog.Error(err.Error())
				return utilities.Response(c, &utilities.ResponseRequest{
					Code:  http.StatusInternalServerError,
					Error: utilities.ErrorRequest(err, http.StatusInternalServerError),
				})
			}

			//If request method POST | PUT | DELETE will print reqData
			if c.Request().Method == http.MethodPost ||
				c.Request().Method == http.MethodPut ||
				c.Request().Method == http.MethodDelete {
				requestLog.Info(string(reqData))
			}

			c.Request().Body = ioutil.NopCloser(bytes.NewReader(reqData))
			return next(c)
		}
	})
	e.Use(middleware.BodyLimit("10M"))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			"Authorization",
			"Version",
		},
	}))
	e.Use(middleware.RateLimiterWithConfig(rateLimitConfig()))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:   middleware.DefaultSkipper,
		StackSize: 1 << 10,
		LogLevel:  log.ERROR,
	}))

	return &Router{e, userRepo}
}

func rateLimitConfig() middleware.RateLimiterConfig {
	return middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 5, Burst: 15, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return utilities.Response(context, &utilities.ResponseRequest{
				Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusTooManyRequests),
			})
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return utilities.Response(context, &utilities.ResponseRequest{
				Error: utilities.ErrorRequest(errors.New(static.ToManyRequest), http.StatusTooManyRequests),
			})
		},
	}
}
