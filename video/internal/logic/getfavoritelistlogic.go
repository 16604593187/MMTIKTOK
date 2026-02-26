package logic

import (
	"context"

	"mini-tiktok/video/internal/svc"
	"mini-tiktok/video/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteListLogic {
	return &GetFavoriteListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFavoriteListLogic) GetFavoriteList(in *video.FavoriteListReq) (*video.FavoriteListResp, error) {
	// todo: add your logic here and delete this line

	return &video.FavoriteListResp{}, nil
}
