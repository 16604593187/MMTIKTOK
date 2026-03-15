// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"

	"mini-tiktok/api/internal/svc"
	"mini-tiktok/api/internal/types"
	"mini-tiktok/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshReq) (resp *types.RefreshResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.Refresh(l.ctx, &user.RefreshReq{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return nil, err
	}

	if rpcResp.StatusCode != 0 {
		return &types.RefreshResp{
			Response: types.Response{
				StatusCode: fmt.Sprintf("%d", rpcResp.StatusCode),
				StatusMsg:  rpcResp.StatusMsg,
			},
		}, nil
	}

	// 业务完全成功
	return &types.RefreshResp{
		Response: types.Response{
			StatusCode: "0",
			StatusMsg:  rpcResp.StatusMsg,
		},
		Token:        rpcResp.AccessToken,
		RefreshToken: rpcResp.RefreshToken,
	}, nil
}
