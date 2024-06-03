package router

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"pcstakehometest/config"
	"pcstakehometest/model"
	"pcstakehometest/package/jwt"
	"pcstakehometest/static"
	"pcstakehometest/utilities"
)

func (r *Router) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		AuthorizationHeader := c.Request().Header.Get("Authorization")
		Authorization := strings.Split(AuthorizationHeader, " ")
		if len(Authorization) > 1 {
			result, err := jwt.ParseClaim(Authorization[1], config.Get().Auth.Secret)
			if err != nil {
				r.Logger.Error(err.Error())
				return utilities.Response(c, &utilities.ResponseRequest{
					Code:  http.StatusUnauthorized,
					Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
				})
			}

			ctx := c.Request().Context()

			// Check exist user
			userDetail, err := r.userRepo.Find(ctx, &model.Users{
				ID: result.Data.UserID,
			})
			if err != nil {
				r.Logger.Error(err.Error())
				return utilities.Response(c, &utilities.ResponseRequest{
					Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
				})
			}

			ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwt.InternalClaimData{
				UserID: userDetail.ID,
				Role:   userDetail.Role,
			})

			c.SetRequest(c.Request().WithContext(ctx))

		} else {
			return utilities.Response(c, &utilities.ResponseRequest{
				Code:  http.StatusUnauthorized,
				Error: utilities.ErrorRequest(errors.New(static.Authorization), http.StatusUnauthorized),
			})
		}
		return next(c)
	}
}
