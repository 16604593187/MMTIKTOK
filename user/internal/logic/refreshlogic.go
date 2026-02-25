package logic

import (
	"context"
	"errors"
	"time"

	"mini-tiktok/user/internal/svc"
	"mini-tiktok/user/user"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshLogic {
	return &RefreshLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefreshLogic) Refresh(in *user.RefreshReq) (*user.RefreshResp, error) {
	secretKey := l.svcCtx.Config.JwtAuth.AccessSecret
	//Parse函数校验签名与过期时间
	token, err := jwt.Parse(in.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		// 如果篡改了签名，或者 RefreshToken 过期，拒绝
		return nil, errors.New("invalid or expired refresh token")
	}
	//提取UserId
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token format")
	}
	//JWT解析出的字段数据类型为float64，需要转换
	userIdFloat, ok := claims["userId"].(float64)
	if !ok {
		return nil, errors.New("userId not found in token")
	}
	userId := uint64(userIdFloat)
	//与注册和登录相同的
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	refreshExpire := l.svcCtx.Config.JwtAuth.RefreshExpire
	newAccessToken, newRefreshToken, err := getJwtTokens(secretKey, now, accessExpire, refreshExpire, userId)

	return &user.RefreshResp{
		StatusCode: 0,
		StatusMsg: "refresh success",
		AccessToken: newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
