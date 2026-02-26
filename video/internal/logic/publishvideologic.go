package logic

import (
	"context"
	"time"

	"mini-tiktok/video/internal/svc"
	"mini-tiktok/video/video" 
	"mini-tiktok/video/model" 

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishVideoLogic {
	return &PublishVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishVideoLogic) PublishVideo(in *video.PublishVideoReq) (*video.PublishVideoResp, error) {
	newVideo := &model.Video{
		//mini-tiktok-mysql.sql 里设置了 AUTO_INCREMENT，MySQL 会自动生成主键
		UserId:     uint64(in.UserId),
		Title:      in.Title,
		PlayUrl:    in.PlayUrl,
		CoverUrl:   in.CoverUrl, // V1.0 先传默认图或空字符串
		CreateTime: time.Now().Unix(),
	}

	//调用 GORM 写入 MySQL
	err := l.svcCtx.Db.Create(newVideo).Error
	if err != nil {
		logx.Errorf("视频写入数据库失败: %v", err)
		return nil, err
	}

	logx.Infof("视频信息落盘成功！VideoID: %d, 标题: %s", newVideo.Id, newVideo.Title)

	// 3. 返回成功响应
	return &video.PublishVideoResp{
		StatusCode: 0, 
		StatusMsg:  "视频发布成功",
	}, nil

}
