package enum

import (
	"fmt"

	"pcstakehometest/static"
)

type RoleType int

const (
	RoleTypeSeller RoleType = 1
	RoleTypeBuyer  RoleType = 2
)

func (t RoleType) String() string {
	switch t {
	case RoleTypeSeller:
		return "Seller"
	case RoleTypeBuyer:
		return "Buyer"
	default:
		return "Unknown"
	}
}

func (t RoleType) IsValid() error {
	switch t {
	case RoleTypeSeller, RoleTypeBuyer:
		return nil
	}
	return fmt.Errorf(static.DataNotFound, "Role")
}
