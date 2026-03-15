package logic

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"mini-tiktok/video/internal/svc"
	"mini-tiktok/video/model"
	"mini-tiktok/video/model/KafkaMessage"
	"mini-tiktok/video/video"

	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavoriteActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteActionLogic {
	return &FavoriteActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FavoriteActionLogic) FavoriteAction(in *video.FavoriteReq) (*video.FavoriteResp, error) {
	//提取请求参数
	videoId, err := strconv.ParseUint(in.VideoId, 10, 64)
	if err != nil {
		return nil, errors.New("无效的视频ID")
	}

	userid, err := strconv.ParseUint(in.UserId, 10, 64)
	if err != nil {
		return nil, errors.New("无效的用户ID")
	}

	actionTypeStr := in.ActionType
	actionTypeInt, err := strconv.Atoi(actionTypeStr)

	conn := l.svcCtx.Redis.NewRedisConn()
	defer conn.Close()

	//调用Lua接口
	changed, err := l.svcCtx.Redis.ExecFavoriteAction(conn, int64(userid), int64(videoId), actionTypeInt)
	if err != nil {
		l.Logger.Errorf("Redis Lua 脚本执行失败: %v", err)
		return nil, err
	}

	//如果状态未改变，直接阻断
	if !changed {
		l.Logger.Infof("拦截到重复操作：用户 %d 对视频 %d 尝试重复执行动作 %s，已直接放行", userid, videoId, actionTypeStr)
		return &video.FavoriteResp{
			StatusCode: STATUS_SUCCESS,
			StatusMsg:  STATUS_SUCCESS_MSG,
		}, nil
	}
	favorite := model.Favorite{
		UserId:  userid,
		VideoId: videoId,
	}
	var op string
	switch actionTypeStr {
	case FAVORITE_UPDATE:
		favorite.CreateTime = time.Now().Unix()
		op = OP_INSERT
	case FAVORITE_DELETE:
		op = OP_DELETE
	default:
		return nil, errors.New(STATUS_FAIL_PARAM_MSG)
	}
	//打包并投递 Kafka 消息
	cols, err := json.Marshal(favorite)
	if err != nil {
		return nil, err
	}

	msg := KafkaMessage.MsgInfo{
		Op:      op,
		Model:   MODEL_FAVORITE,
		Columns: string(cols),
	}

	marshalStr, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	kafkaMsg := kafka.Message{
		Value: marshalStr,
	}
	err = l.svcCtx.FavouriteBatcher.Push(kafkaMsg)
	if err != nil {
		l.Logger.Errorf("投递 Kafka 失败，准备回滚 Redis 缓存: %v", err)

		// 补偿机制：直接删除相关的 Redis Key
		userCacheKey := model.Favorite{}.CacheKey(userid)
		videoCountKey := model.Favorite{}.CountCacheKey(videoId)

		_, delErr := l.svcCtx.Redis.Del(userCacheKey, videoCountKey)
		if delErr != nil {
			l.Logger.Errorf("严重错误：Kafka 写入失败且 Redis 回滚也失败，发生脏数据死锁！userId: %d, videoId: %d, err: %v", userid, videoId, delErr)
		}
		return nil, errors.New("请稍后再试")
	}
	return &video.FavoriteResp{
		StatusCode: STATUS_SUCCESS,
		StatusMsg:  STATUS_SUCCESS_MSG,
	}, nil
}
