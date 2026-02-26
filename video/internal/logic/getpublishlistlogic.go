package logic

import (
	"context"

	"mini-tiktok/video/internal/svc"
	"mini-tiktok/video/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPublishListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPublishListLogic {
	return &GetPublishListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPublishListLogic) GetPublishList(in *video.PublishListReq) (*video.PublishListResp, error) {
	// todo: add your logic here and delete this line

	return &video.PublishListResp{}, nil
}
