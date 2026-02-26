package logic

import (
	"context"

	"mini-tiktok/video/internal/svc"
	"mini-tiktok/video/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavoriteActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteActionLogic {
	return &FavoriteActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FavoriteActionLogic) FavoriteAction(in *video.FavoriteReq) (*video.FavoriteResp, error) {
	// todo: add your logic here and delete this line

	return &video.FavoriteResp{}, nil
}
