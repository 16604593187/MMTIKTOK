package logic

import (
	"context"

	"mini-tiktok/user/internal/svc"
	"mini-tiktok/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowListLogic) FollowList(in *user.FollowListReq) (*user.FollowListResp, error) {
	// todo: add your logic here and delete this line

	return &user.FollowListResp{}, nil
}
