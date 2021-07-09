package model


import (
	"examresult/config"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// DB 数据库连接单例
var DB *gorm.DB

// Init 初始化 MySQL 连接。
// 连接失败将导致程序 Fatal 退出。
// 连接的具体方法参考： https://gorm.io/docs/connecting_to_the_database.html#MySQL
func Init() {
	log.Info("init database")

	// Data Source Name: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	// Refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details.
	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Database.Username,
		config.Database.Password,
		config.Database.Protocol,
		config.Database.Address,
		config.Database.DBName,
	)

	logger := log.WithFields(log.Fields{"dsn": dsn})

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.WithField("err", err).Fatal("failed to connect database")
		return
	}

	// 设置连接池：https://gorm.io/docs/connecting_to_the_database.html#Connection-Pool
	sqlDB, err := db.DB()
	if err != nil {
		logger.WithField("err", err).Fatal("failed to set Connection Pool")
		return
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(25)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err = migrate(db); err != nil {
		logger.WithField("err", err).Fatal("failed to set migrate database")
		return
	}

	DB = db
}