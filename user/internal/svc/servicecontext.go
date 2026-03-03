package svc

import (
	"log"
	"mini-tiktok/user/internal/config"
	"mini-tiktok/user/model"
	"mini-tiktok/user/model/redisCache"

	"gorm.io/gorm"
	"github.com/ncghost1/snowflake-go"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redisCache.RedisPool
	Db     *gorm.DB
	SnowflakeNode *snowflake.SnowFlake
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := model.InitGorm(c.DbConfig)
	if err != nil {
		log.Fatalln(err)
		return nil
	} //
	pool := redisCache.NewRedisPool(c)
	conn := pool.NewRedisConn()
	_, err = conn.Do("PING")
	defer conn.Close()
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	sf, err := snowflake.New(c.WorkerId)
	if err != nil {
		panic("初始化雪花算法失败: " + err.Error())
	}
	return &ServiceContext{
		Config: c,
		Redis:  pool,
		Db:     db,
		SnowflakeNode: sf,
	}
}
