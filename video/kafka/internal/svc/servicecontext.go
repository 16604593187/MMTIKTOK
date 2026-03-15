package svc

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	
	"mini-tiktok/video/kafka/internal/config"
	// 如果需要 Redis 也可以引入你的 redisCache 包
)

type ServiceContext struct {
	Config config.Config
	Db     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	dsn := c.DbConfig.Username + ":" + c.DbConfig.Password + "@tcp(" + c.DbConfig.Path + ":3306)/" + c.DbConfig.Dbname + "?" + c.DbConfig.Config
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Kafka消费者连接MySQL失败:", err)
		return nil
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(c.DbConfig.MaxIdleConns) 
		sqlDB.SetMaxOpenConns(c.DbConfig.MaxOpenConns) 
	}

	return &ServiceContext{
		Config: c,
		Db:     db,
	}
}