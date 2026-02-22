package logic

import (
	"context"

	"mini-tiktok/user/internal/svc"
	"mini-tiktok/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowActionLogic {
	return &FollowActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowActionLogic) FollowAction(in *user.FollowActionReq) (*user.FollowActionResp, error) {
	// todo: add your logic here and delete this line

	return &user.FollowActionResp{}, nil
}
