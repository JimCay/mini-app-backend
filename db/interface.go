package db

import (
	"context"
	"tg-backend/db/model"
)

type Storage interface {
	UserStorage

	PointStorage

	TaskStorage
}

type UserStorage interface {
	InsertUser(ctx context.Context, user *model.User, inviteCode uint64) error

	GetUser(ctx context.Context, id uint64) (*model.User, error)

	UpdateUserDays(ctx context.Context, user *model.User) error

	GetFriends(ctx context.Context, id uint64) ([]model.MyInvitee, error)
}

type PointStorage interface {
	GetPoint(ctx context.Context, id uint64) (*model.Point, error)

	UpdatePoint(ctx context.Context, point *model.Point) error

	Ranks(ctx context.Context) ([]model.Rank, error)

	MyRank(ctx context.Context, id uint64) (int64, error)
}

type TaskStorage interface {
	GetTasks(ctx context.Context) ([]model.Task, error)

	GetTask(ctx context.Context, id uint64) (*model.Task, error)

	GetUserStatus(ctx context.Context, taskId, userId uint64) (bool, error)

	UserTask(ctx context.Context, userId uint64) ([]model.UserTask, error)

	TaskDone(ctx context.Context, task *model.Task, userId uint64) error
}
