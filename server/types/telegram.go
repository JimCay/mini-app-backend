package types

import (
	"tg-backend/db/model"
	"time"
)

type TelegramUser struct {
	ID              int    `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Username        string `json:"username"`
	LanguageCode    string `json:"language_code"`
	IsPremium       bool   `json:"is_premium"`
	AllowsWriteToPm bool   `json:"allows_write_to_pm"`
}

func (tu *TelegramUser) GetUser() *model.User {
	return &model.User{
		Id:        uint64(tu.ID),
		FirstName: tu.FirstName,
		LastName:  tu.LastName,
		Username:  tu.Username,
		CreatedAt: time.Now(),
	}
}
