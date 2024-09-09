package service

import (
	"tg-backend/config"
	"tg-backend/db"
)

type ServiceManager struct {
	User  *UserService
	Point *PointService
	Task  *TaskService
}

func NewServiceManager(storage db.Storage, config *config.Config) *ServiceManager {
	user := NewUserService(storage)
	point := NewPointService(storage)
	task := NewTaskService(storage)

	return &ServiceManager{user, point, task}
}
