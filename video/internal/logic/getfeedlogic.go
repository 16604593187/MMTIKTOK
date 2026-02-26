package logic

import (
	"context"
	"time"

	"mini-tiktok/video/internal/svc" 
	"mini-tiktok/video/model"      
	"mini-tiktok/video/video"      

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFeedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedLogic {
	return &GetFeedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFeedLogic) GetFeed(in *video.FeedReq) (*video.FeedResp, error) {
	//处理前端传来的时间戳
	// 如果前端传 0，取当前系统时间
	latestTime := in.LatestTime
	if latestTime == 0 {
		latestTime = time.Now().Unix()
	}
	//使用 GORM 从 MySQL 拉取视频 (每次拉取 30 条)
	var dbVideos []*model.Video
	err := l.svcCtx.Db.Where("create_time < ?", latestTime).
		Order("create_time desc").
		Limit(30).
		Find(&dbVideos).Error

	if err != nil {
		logx.Errorf("拉取 Feed 流失败: %v", err)
		return nil, err
	}

	//数据库模型转化为RPC通信模型
	var rpcVideoList []*video.Video
	for _, v := range dbVideos {
		rpcVideoList = append(rpcVideoList, &video.Video{
			ID:            v.Id,
			PlayURL:       v.PlayUrl,
			CoverURL:      v.CoverUrl,
			Title:         v.Title,
			FavoriteCount: 0, // V1.0 先写死
			CommentCount:  0, 
			IsFavorite:    false,
		})
	}
	var nextTime int64
	if len(dbVideos) > 0 {
		nextTime = dbVideos[len(dbVideos)-1].CreateTime
	} else {
		nextTime = time.Now().Unix() // 没拉到数据，重置为当前时间
	}

	return &video.FeedResp{
		StatusCode: "0",
		StatusMsg:  "获取 Feed 流成功",
		VideoList:  rpcVideoList,
		NextTime:   nextTime,
	}, nil
}
