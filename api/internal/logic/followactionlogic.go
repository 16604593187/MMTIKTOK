// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"mini-tiktok/api/internal/svc"
	"mini-tiktok/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowActionLogic {
	return &FollowActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowActionLogic) FollowAction(req *types.FollowActionReq) (resp *types.FollowActionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
