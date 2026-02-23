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

type RegisterUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterUserLogic {
	return &RegisterUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterUserLogic) RegisterUser(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	//调用 User RPC 的 Register 方法
	rpcResp, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterReq{
		Username: req.Username,
		Password: req.Password,
	})

	// 如果 RPC 调用失败
	if err != nil {
		return nil, err
	}

	//将 RPC 的返回值，包装成 API 网关需要的 HTTP JSON 格式返回给前端
	return &types.RegisterResp{
		Response: types.Response{
			StatusCode: rpcResp.StatusCode,
			StatusMsg:  rpcResp.StatusMsg,
		},
		UserID: rpcResp.UserID,
		// Token: rpcResp.Token,
	}, nil

}
