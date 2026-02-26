// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
	"strconv"
	"mini-tiktok/api/internal/svc"
	"mini-tiktok/api/internal/types"
	"mini-tiktok/video/video"
	"github.com/zeromicro/go-zero/core/logx"
)

type PublishActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishActionLogic {
	return &PublishActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishActionLogic) PublishAction(req *types.PublishReq, handler *multipart.FileHeader, file multipart.File) (resp *types.PublishResp, err error) {

	var uid string
	if userIdVal := l.ctx.Value("userId"); userIdVal != nil {
		uid = fmt.Sprintf("%v", userIdVal)
	} else {
		uid = "unknown_user"
	}
	// 获取当前运行目录下的 public/videos
	videoDir := filepath.Join("public", "videos")
	if err := os.MkdirAll(videoDir, os.ModePerm); err != nil {
		logx.Errorf("自动创建目录失败: %v", err)
		return nil, err
	}

	// 拼接文件名
	filename := fmt.Sprintf("%s_%d_%s", uid, time.Now().Unix(), handler.Filename)
	savePath := filepath.Join(videoDir, filename)

	//创建文件并写入
	dst, err := os.Create(savePath)
	if err != nil {
		logx.Errorf("创建本地文件失败: %v", err)
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		logx.Errorf("保存视频流失败: %v", err)
		return nil, err
	}

	logx.Infof("视频保存在本地成功: %s", savePath)

	playURL := fmt.Sprintf("http://127.0.0.1:8888/videos/%s", filename)
	logx.Infof("准备写入数据库的 PlayURL: %s", playURL)

	uidInt64, _ := strconv.ParseInt(uid, 10, 64)

	rpcResp, err := l.svcCtx.VideoRpc.PublishVideo(l.ctx, &video.PublishVideoReq{
		UserId:   uidInt64,
		Title:    req.Title,
		PlayUrl:  playURL,
		CoverUrl: "http://127.0.0.1:8888/default_cover.jpg", // 先写死一个默认封面
	})

	if err != nil {
		logx.Errorf("调用 Video RPC 失败: %v", err)
		return nil, err
	}

	return &types.PublishResp{
		Response: types.Response{
			StatusCode: strconv.FormatInt(rpcResp.StatusCode,10),
			StatusMsg:  rpcResp.StatusMsg,
		},
	}, nil
}
