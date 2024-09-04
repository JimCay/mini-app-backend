package types

import (
	"tg-backend/db/model"
	"time"
)

type LoginInfo struct {
	InviteCode string `json:"invite_code"`
}

type Invite struct {
	Code string `json:"code"`
}

type Point struct {
	Id    uint64 `json:"id"`
	Value uint64 `json:"value"`
	Rate  uint64 `json:"rate"`
}

func GetPoint(point *model.Point) *Point {
	return &Point{
		Id:    point.Id,
		Value: point.Value,
		Rate:  point.Rate,
	}
}

func (p *Point) GetPoint(id uint64) *model.Point {
	return &model.Point{
		Id:        id,
		Value:     p.Value,
		Rate:      p.Rate,
		UpdatedAt: time.Now(),
	}
}

type Friend struct {
	Invitor uint64 `json:"invitor"`
	Invitee uint64 `json:"invitee"`
	Reward  uint64 `json:"reward"`
}

func GetFriend(friend model.Friend) Friend {
	return Friend{
		Invitor: friend.Invitor,
		Invitee: friend.Invitee,
		Reward:  friend.Value,
	}
}
