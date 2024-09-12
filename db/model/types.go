package model

import (
	"tg-backend/config"
	"time"
)

const TaskTypeFriend int = 1

const TaskTypeX int = 2

const DELETE int8 = 0
const NORMAL int8 = 1

const EnergyOfTime = float32(3600*24) / float32(config.DefaultDayLimit)

type User struct {
	Id        uint64 `gorm:"primary_key" `
	FirstName string
	LastName  string
	UserName  string
	UpdatedAt time.Time
	CreatedAt time.Time
	LoginDays int32
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

func (p *Point) GetEnergy() int32 {
	gap := time.Now().Unix() - p.UpdatedAt.Unix()
	gapEnergy := int32(0)
	if gap > 0 {
		gapEnergy = int32(float32(gap) / EnergyOfTime)
	}
	energy := p.Energy + gapEnergy
	if energy > p.Limit {
		energy = p.Limit
	}
	return energy
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
