package dto

import (
	"fmt"

	"pcstakehometest/static"
)

type LoginRequest struct {
	Username string
	Password string
}

func (d *LoginRequest) Validate() error {
	if d.Username == "" {
		return fmt.Errorf(static.EmptyValue, "email")
	}
	if d.Password == "" {
		return fmt.Errorf(static.EmptyValue, "password")
	}
	return nil
}

type Response struct {
	Token        string
	RefreshToken string
}
