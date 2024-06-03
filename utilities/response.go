package utilities

import (
	"github.com/labstack/echo/v4"
)

type ResponseRequest struct {
	Data   interface{}
	Code   int
	Status string
	Error  error
}

// Response :
func Response(c echo.Context, r *ResponseRequest) error {
	resp := make(map[string]interface{})
	if r.Error != nil {
		errInfo, ok := r.Error.(*Error)
		if ok {
			resp["Code"] = errInfo.StatusCode
			resp["Status"] = errInfo.Error()
			return c.JSON(errInfo.StatusCode, resp)
		}
	}
	resp["Code"] = r.Code
	resp["Status"] = r.Status
	if r.Data != nil {
		resp["Data"] = r.Data
	}
	return c.JSON(r.Code, resp)
}
