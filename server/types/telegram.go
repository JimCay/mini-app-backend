package types

import (
	"tg-backend/db/model"
	"time"
)

type TelegramUser struct {
	ID              uint64 `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	UserName        string `json:"username"`
	LanguageCode    string `json:"language_code"`
	IsPremium       bool   `json:"is_premium"`
	AllowsWriteToPm bool   `json:"allows_write_to_pm"`
}

func (tu *TelegramUser) GetUser() *model.User {
	return &model.User{
		Id:        tu.ID,
		FirstName: tu.FirstName,
		LastName:  tu.LastName,
		UserName:  tu.UserName,
		CreatedAt: time.Now(),
	}
}

func (tu *TelegramUser) UserPoint(value uint64, pos int32) UserPoint {
	return UserPoint{
		Id:       tu.ID,
		Value:    value,
		Pos:      pos,
		Name:     tu.FirstName + " " + tu.LastName,
		UserName: tu.UserName,
	}
}
