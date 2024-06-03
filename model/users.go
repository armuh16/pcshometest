package model

import (
	"pcstakehometest/enum"
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID        int
	Username  string
	Password  string        `json:"-"`
	Role      enum.RoleType `json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `json:"-"`
}
