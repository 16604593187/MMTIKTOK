package logic

import (
	"context"

	"mini-tiktok/video/internal/svc"
	"mini-tiktok/video/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentListLogic {
	return &GetCommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentListLogic) GetCommentList(in *video.CommentListReq) (*video.CommentListResp, error) {
	// todo: add your logic here and delete this line

	return &video.CommentListResp{}, nil
}
