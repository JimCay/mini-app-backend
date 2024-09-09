package db

import (
	"errors"
	"strings"
	"tg-backend/config"
	"tg-backend/db/mysql"
)

const MYSQL = "MYSQL"

func Setup(config config.DbConfig) (Storage, error) {
	if strings.ToUpper(config.Type) == MYSQL {
		db, err := mysql.InitMysql(config)
		if err != nil {
			return nil, err
		}
		return db, nil
	}
	return nil, errors.New("db type not supported")
}
