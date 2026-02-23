// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"mini-tiktok/api/internal/svc"
	"mini-tiktok/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginUserLogic {
	return &LoginUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginUserLogic) LoginUser(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line

	return
}
