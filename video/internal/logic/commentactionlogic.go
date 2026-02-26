package logic

import (
	"context"

	"mini-tiktok/video/internal/svc"
	"mini-tiktok/video/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentActionLogic {
	return &CommentActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommentActionLogic) CommentAction(in *video.CommentReq) (*video.CommentResp, error) {
	// todo: add your logic here and delete this line

	return &video.CommentResp{}, nil
}
