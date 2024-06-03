package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"pcstakehometest/enum"
	"time"

	"gorm.io/gorm"
	"pcstakehometest/static"
)

type Transactions struct {
	ID          int
	BuyerID     int `json:"-"`
	SellerID    int `json:"-"`
	Origin      string
	Destination string
	Items       ItemsTransaction
	GrandTotal  float64
	Status      enum.TransactionStatusType `json:"-"`
	Coupons     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `json:"-"`

	// Relations
	Seller *Users `json:",omitempty" gorm:"<-:false;foreignKey:SellerID;references:ID;"`
	Buyer  *Users `json:",omitempty" gorm:"<-:false;foreignKey:BuyerID;references:ID;"`

	// Attribute
	StatusTransaction string `gorm:"<-:false;-;"`
}

type ItemsTransaction []ProductTransaction

type ProductTransaction struct {
	ID          int
	Name        string
	Description string
	Price       float64
}

func (j ItemsTransaction) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *ItemsTransaction) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf(static.SomethingWrong)
	}

	result := ItemsTransaction{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	*j = result

	return nil
}
