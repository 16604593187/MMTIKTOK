// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"strconv"

	"mini-tiktok/api/internal/svc"
	"mini-tiktok/api/internal/types"
	"mini-tiktok/video/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type FeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {
	return &FeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FeedLogic) Feed(req *types.FeedReq) (resp *types.FeedResp, err error) {
	//将前端传来的string类型时间戳转化为int64
	var latestTime int64
	if req.LatestTime != nil && *req.LatestTime != "" {
		latestTime, _ = strconv.ParseInt(*req.LatestTime, 10, 64)
	}
	//调用Video RPC拉取视频列表
	rpcResp, err := l.svcCtx.VideoRpc.GetFeed(l.ctx, &video.FeedReq{
		LatestTime: latestTime,
	})
	if err != nil {
		logx.Errorf("调用 VideoRpc GetFeed 失败: %v", err)
		return nil, err
	}
	//转换列表为网关要求的types video形式
	var apiVideoList []types.Video
	for _, v := range rpcResp.VideoList {
		apiVideoList = append(apiVideoList, types.Video{
			ID:            v.ID,
			PlayURL:       v.PlayURL,
			CoverURL:      v.CoverURL,
			Title:         v.Title,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Author: types.User{
				ID:   1,
				Name: "MVP_User",
			},
		})
	}
	return &types.FeedResp{
		Response: types.Response{
			StatusCode: rpcResp.StatusCode,
			StatusMsg:  rpcResp.StatusMsg,
		},
		VideoList: apiVideoList,
		NextTime:  &rpcResp.NextTime,
	}, nil
}
