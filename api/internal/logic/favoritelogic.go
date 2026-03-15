// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"
	"mini-tiktok/api/internal/svc"
	"mini-tiktok/api/internal/types"
	"mini-tiktok/video/videorpc"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteLogic {
	return &FavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteLogic) Favorite(req *types.FavoriteReq) (resp *types.FavoriteResp, err error) {
	uid, _ := l.ctx.Value("userId").(json.Number).Int64()

	// 2. 调用 Video RPC 服务
	// 这个接口会触发 video/internal/logic/favoriteactionlogic.go 里的逻辑
	_, err = l.svcCtx.VideoRpc.FavoriteAction(l.ctx, &videorpc.FavoriteReq{
		UserId:     strconv.FormatUint(uint64(uid), 10),
		VideoId:    req.VideoId,    // 确保 types.go 中 VideoId 已正确定义
		ActionType: req.ActionType, // 1-点赞，2-取消点赞
	})

	if err != nil {
		l.Errorf("调用 VideoRpc.FavoriteAction 失败: %v", err)
		return nil, err
	}

	return &types.FavoriteResp{
		Response: types.Response{
			StatusCode: "0",
			StatusMsg:  "success",
		},
	}, nil

}
