// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"mini-tiktok/api/internal/svc"
	"mini-tiktok/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPublishListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPublishListLogic {
	return &GetPublishListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPublishListLogic) GetPublishList(req *types.PublishListReq) (resp *types.PublishListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
