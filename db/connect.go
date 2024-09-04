package db

import (
	"errors"
	"gorm.io/gorm"
	"strings"
	"tg-backend/config"
	"tg-backend/db/mysql"
)

const MYSQL = "MYSQL"

type Storage struct {
	db *gorm.DB
}

func Setup(config config.DbConfig) (*Storage, error) {
	if strings.ToUpper(config.Type) == MYSQL {
		db, err := mysql.InitMysql(config)
		if err != nil {
			return nil, err
		}
		return &Storage{db}, nil
	}
	return nil, errors.New("db type not supported")
}
