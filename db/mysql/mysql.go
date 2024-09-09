package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"tg-backend/config"
	"tg-backend/db/model"
	"tg-backend/pkg/log"
)

type MysqlStorage struct {
	db *gorm.DB
}

func InitMysql(config config.DbConfig) (*MysqlStorage, error) {
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
	dbUri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Asia%%2fShanghai&timeout=30s",
		config.User,
		config.Password,
		config.IP,
		config.Port,
		config.Name)

	dialector := mysql.New(mysql.Config{
		DSN:                       dbUri, // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	})
	mysqlDb, err := gorm.Open(dialector, gormConfig)
	if nil != err {
		log.Error("db Setup err: %v", err)
		return nil, err
	}
	if config.NumberShard > 1 {
		//middleware := sharding.Register(sharding.Config{
		//	ShardingKey:         "file_id",
		//	NumberOfShards:      uint(config.NumberShard),
		//	PrimaryKeyGenerator: sharding.PKSnowflake,
		//}, "xxxxx")
		//MysqlDb.Use(middleware)
	}

	if err = Migrator(mysqlDb); err != nil {
		return nil, err
	}
	log.Info("db init complete")
	return &MysqlStorage{mysqlDb}, nil
}

func Migrator(db *gorm.DB) error {
	if err := db.Migrator().
		AutoMigrate(
			&model.User{},
			&model.Friend{},
			&model.Point{},
			&model.Task{},
			&model.UserTask{}); err != nil {
		return err
	}
	return nil
}
