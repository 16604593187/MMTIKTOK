// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"mini-tiktok/api/internal/svc"
	"mini-tiktok/api/internal/types"
	"mini-tiktok/user/user"

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
	rpcResp, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		Response: types.Response{
			StatusCode: rpcResp.StatusCode,
			StatusMsg:  rpcResp.StatusMsg,
		},
		UserID:       rpcResp.UserID,
		Token:        &rpcResp.AccessToken,
		RefreshToken: rpcResp.RefreshToken,
	}, nil
}
