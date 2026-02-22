package svc

import (
	"log"
	"mini-tiktok/user/internal/config"
	"mini-tiktok/user/model"
	"mini-tiktok/user/model/redisCache"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redisCache.RedisPool
	Db     *gorm.DB
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
	return &ServiceContext{
		Config: c,
		Redis:  pool,
		Db:     db,
	}
}
