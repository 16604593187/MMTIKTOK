package main

import (
	"flag"

	"github.com/zeromicro/go-zero/core/conf"
	
	"mini-tiktok/video/kafka/internal/config"
	"mini-tiktok/video/kafka/internal/logic"
	"mini-tiktok/video/kafka/internal/svc"
)

var configFile = flag.String("f", "etc/video-kafka.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 初始化消费端服务上下文
	svcCtx := svc.NewServiceContext(c)

	consumer := logic.NewFavoriteConsumer(svcCtx)
	consumer.Start() 
}