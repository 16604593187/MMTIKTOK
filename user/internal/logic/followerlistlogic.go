package logic

import (
	"context"

	"mini-tiktok/user/internal/svc"
	"mini-tiktok/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowerListLogic {
	return &FollowerListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowerListLogic) FollowerList(in *user.FollowerListReq) (*user.FollowerListResp, error) {
	// todo: add your logic here and delete this line

	return &user.FollowerListResp{}, nil
}
