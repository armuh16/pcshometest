package model

import (
	"time"

	"gorm.io/gorm"
)

type Products struct {
	ID          int
	Name        string
	Description string
	Price       float64
	SellerID    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `json:"-"`

	// Relations
	Seller *Users `json:",omitempty" gorm:"<-:false;foreignKey:SellerID;references:ID;"`
}
