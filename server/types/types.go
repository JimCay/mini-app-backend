package types

import (
	"tg-backend/db/model"
	"time"
)

type LoginInfo struct {
	InviteCode string `json:"invite_code"`
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
