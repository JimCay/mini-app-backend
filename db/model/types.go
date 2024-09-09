package model

import (
	"time"
)

const TASK_TYPE_FRIEND int = 1

const DELETE int8 = 0
const NORMAL int8 = 1

type User struct {
	Id        uint64 `gorm:"primary_key" `
	FirstName string
	LastName  string
	UserName  string
	CreatedAt time.Time
}

type Friend struct {
	Id           uint64 `gorm:"primary_key" `
	Invitor      uint64 `gorm:"index:idx_invitor"`
	Invitee      uint64 `gorm:"index:idx_invitee"`
	Value        int32
	InviteeValue int32
	CreatedAt    time.Time
}

type Point struct {
	Id        uint64 `gorm:"primary_key" `
	Value     uint64 `gorm:"index:idx_value"`
	Rate      int32
	Limit     int32
	Energy    int32
	UpdatedAt time.Time
}

type Task struct {
	Id          uint64 `gorm:"primary_key" `
	Name        string
	Description string `gorm:"type:text"`
	TaskType    int
	Condition   int
	Reward      int32
	Status      int8
	CreatedAt   time.Time
}

type UserTask struct {
	Id        uint64 `gorm:"primary_key" `
	UserId    uint64 `gorm:"uniqueIndex:idx_user_id;not null"`
	TaskId    uint64 `gorm:"uniqueIndex:idx_user_id;not null"`
	CreatedAt time.Time
}

type MyInvitee struct {
	User
	Reward int32
}

type Rank struct {
	User
	Value uint64
}
