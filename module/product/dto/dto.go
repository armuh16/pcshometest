package dto

import (
	"errors"
	"fmt"

	"pcstakehometest/enum"
	"pcstakehometest/model"
	"pcstakehometest/static"
)

type CreateRequest struct {
	Name        string
	Description string
	Price       float64
	SellerID    int
	RoleID      enum.RoleType
}

func (d *CreateRequest) Validate() error {
	if d.SellerID <= 0 {
		return fmt.Errorf(static.EmptyValue, "SellerID")
	}
	if d.Description == "" {
		return fmt.Errorf(static.EmptyValue, "Description")
	}
	if d.Price <= 0 {
		return fmt.Errorf(static.MinValue, "Price", 0)
	}
	if err := d.RoleID.IsValid(); err != nil {
		return err
	}
	if d.RoleID != enum.RoleTypeSeller {
		return errors.New(static.Authorization)
	}
	return nil
}

type FindAllRequest struct {
	SellerID int
}

func (d *FindAllRequest) Validate() error {
	if d.SellerID <= 0 {
		return fmt.Errorf(static.EmptyValue, "SellerID")
	}
	return nil
}

type FindRequest model.Products
