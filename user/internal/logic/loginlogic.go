package logic

import (
	"context"
	"errors"
	"mini-tiktok/user/internal/logic/utils"
	"mini-tiktok/user/internal/svc"
	"mini-tiktok/user/model"
	"mini-tiktok/user/user"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	//登录逻辑
	var userInfo *model.User
	err := l.svcCtx.Db.Where(&model.User{Username: in.Username}).Take(&userInfo).Error //由于用户名唯一，此处提高性能直接使用Take
	if err != nil {
		if err == gorm.ErrRecordNotFound { //用户名不存在
			return &user.LoginResp{
				StatusCode: STATUS_FAIL,
				StatusMsg:  STATUS_USER_NOTEXIST_MSG,
				UserID:     0,
			}, nil
		}
		return nil, err //其余种类错误
	}
	if ok := utils.BcryptCheck(in.Password, userInfo.Password); !ok {
		return &user.LoginResp{
			StatusCode: STATUS_FAIL,
			StatusMsg:  STATUS_WRONG_PASSWORD_MSG, //正常来讲用户名或密码错误会合成一个错误码，此处做了区分。
			UserID:     0,
		}, nil
	}
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	refreshExpire := l.svcCtx.Config.JwtAuth.RefreshExpire
	secretKey := l.svcCtx.Config.JwtAuth.AccessSecret
	accessToken, refreshToken, jti, err := getJwtTokens(secretKey, now, accessExpire, refreshExpire, userInfo.Id)
	if err != nil {
		return nil, err // 签发失败，抛出系统异常
	}

	redisConn := l.svcCtx.Redis.NewRedisConn()
	defer redisConn.Close()

	// 检查账号是否被物理封禁
	banned, err := l.svcCtx.Redis.IsUserBanned(redisConn, userInfo.Id)
	if err == nil && banned {
		return nil, errors.New("该账号涉嫌严重违规，已被冻结") 
	}

	// 将新设备加入zset，触发 FIFO 多端挤占淘汰
	maxDevices := 3 // 限制最多允许 3 台设备同时在线
	tokenTtl := int(refreshExpire) 
	err = l.svcCtx.Redis.AddActiveTokenWithLimit(redisConn, userInfo.Id, jti, maxDevices, tokenTtl)
	if err != nil {
		l.Logger.Errorf("Redis 记录设备在线状态失败: %v", err)
		// 降级处理：即使 Redis 抖动，也允许用户登录
	}
	
	return &user.LoginResp{
		StatusCode:   STATUS_SUCCESS,
		StatusMsg:    STATUS_SUCCESS_MSG,
		UserID:       userInfo.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}


