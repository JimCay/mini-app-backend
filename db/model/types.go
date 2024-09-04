package model

import (
	"time"
)

type User struct {
	Id        uint64 `gorm:"primary_key" `
	FirstName string
	LastName  string
	Username  string
	CreatedAt time.Time
}

type Friends struct {
	Id        uint64 `gorm:"primary_key" `
	Invitor   uint64 `gorm:"index:idx_invitor"`
	Invitee   uint64
	Value     uint64
	CreatedAt time.Time
}

type Point struct {
	Id        uint64 `gorm:"primary_key" `
	Value     uint64
	Rate      uint64
	UpdatedAt time.Time
}
