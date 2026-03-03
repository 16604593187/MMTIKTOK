package logic

import (
	"context"
	"mini-tiktok/user/internal/logic/utils" // 密码加密工具
	"mini-tiktok/user/internal/svc"
	"mini-tiktok/user/model" // 数据库模型
	"mini-tiktok/user/user"  // protobuf 生成的结构体
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	//获取用户名与密码
	username := in.Username
	password := in.Password
	exit := int64(0)
	err := l.svcCtx.Db.Model(&model.User{}).Where(&model.User{Username: username}).Count(&exit).Error
	if err != nil {
		return nil, err
	}
	if exit > 0 { //用户名已存在
		return &user.RegisterResp{
			StatusCode: STATUS_FAIL,
			StatusMsg:  STATUS_USER_EXISTS_MSG,
			UserID:     0,
		}, nil
	}
	//不存在已有用户名则进行后续操作，用雪花算法生成分布式唯一ID
	uuid, err := l.svcCtx.SnowflakeNode.Generate()
	if err != nil {
		return nil, err
	}
	userInfo := &model.User{
		Id:       uuid,
		Username: username,
		Password: utils.BcryptHash(password), //调用hash.go中的函数使用brcypt对密码进行哈希加密
	}
	err = l.svcCtx.Db.Create(&userInfo).Error //创建用户记录
	if err != nil {
		return nil, err
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
	maxDevices := 3
	tokenTtl := int(refreshExpire)
	err = l.svcCtx.Redis.AddActiveTokenWithLimit(redisConn, userInfo.Id, jti, maxDevices, tokenTtl)
	if err != nil {
		l.Logger.Errorf("Redis 记录新注册用户设备状态失败: %v", err)
	}
	return &user.RegisterResp{
		StatusCode:   STATUS_SUCCESS,
		StatusMsg:    STATUS_SUCCESS_MSG,
		UserID:       userInfo.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
