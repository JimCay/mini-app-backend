package service

import (
	"tg-backend/config"
	"tg-backend/db"
)

type ServiceManager struct {
	User  *UserService
	point *PointService
}

func NewServiceManager(storage *db.Storage, config *config.Config) *ServiceManager {
	user := NewUserService(storage)
	point := NewPointService(storage)

	return &ServiceManager{user, point}
}
