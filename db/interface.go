package db

import (
	"context"
	"tg-backend/db/model"
)

type UserStorage interface {
	InsertUser(ctx context.Context, user *model.User, inviteCode uint64) error

	GetUser(ctx context.Context, id uint64) (*model.User, error)

	GetFriends(ctx context.Context, id uint64) ([]model.Friend, error)
}

type PointStorage interface {
	GetPoint(ctx context.Context, id uint64) (*model.Point, error)

	UpdatePoint(ctx context.Context, point *model.Point) error
}
