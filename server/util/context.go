package util

import (
	"context"
	"tg-backend/server/types"
)

const contextAuthKey = "auth"

func NewContext(ctx context.Context, tgUser *types.TelegramUser) context.Context {
	return context.WithValue(ctx, contextAuthKey, tgUser)
}

func FromContext(ctx context.Context) (*types.TelegramUser, bool) {
	auth, ok := ctx.Value(contextAuthKey).(*types.TelegramUser)
	return auth, ok
}
