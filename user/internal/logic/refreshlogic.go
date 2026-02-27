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
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				//自然过期
				l.Logger.Infof("拦截： refreshToken 已自然过期，需重新登录")
				return nil, errors.New("长Token已过期，请重新登录")
			}
		}
		//被篡改或签名错误
		l.Logger.Errorf("拦截：长 Token 签名非法或被篡改: %v", err)
		return nil, errors.New("无效的Token")
	}
	if !token.Valid {
		return nil, errors.New("无效的Token")
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

	redisConn := l.svcCtx.Redis.NewRedisConn()
	defer redisConn.Close()

	banTime, err := l.svcCtx.Redis.GetUserBanTime(redisConn, userId)
	if err == nil && banTime > 0 {
		// 提取这根 Token 的签发时间 iat
		iatFloat, ok := claims["iat"].(float64)
		if ok {
			// 如果 Token 的签发时间 <= 强制下线时间，说明这是案发前的旧 Token
			if int64(iatFloat) <= banTime {
				l.Logger.Errorf("风控拦截：用户 %d 的 Token (签发于 %d) 处于强制下线时间 (%d) 之前，已作废", userId, int64(iatFloat), banTime)
				return nil, errors.New("账号曾存在安全风险，旧凭证已彻底作废，请重新登录")
			}
		}
	}
	//重放拦截
	isUsed, err := l.svcCtx.Redis.CheckTokenBlacklist(redisConn, in.RefreshToken)
	if err == nil && isUsed {
		// 🚨 触发主动反制：记录当前的案发时间戳，斩断此前签发的所有 Token！
		nowTS := time.Now().Unix()
		banTTL := 7 * 24 * 3600 // 保留 7 天的黑名单记录

		errBan := l.svcCtx.Redis.BanUserAtTime(redisConn, userId, nowTS, banTTL)
		if errBan == nil {
			l.Logger.Errorf("🚨 主动反制触发：检测到重放攻击，已将用户 %d 踢下线 (封印时间戳: %d)", userId, nowTS)
		}
		return nil, errors.New("检测到异常请求，为保护账号安全，已强制下线")
	}

	//与注册和登录相同的Token生成逻辑
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	refreshExpire := l.svcCtx.Config.JwtAuth.RefreshExpire
	newAccessToken, newRefreshToken, err := getJwtTokens(secretKey, now, accessExpire, refreshExpire, userId)

	expTimeFloat, ok := claims["exp"].(float64)
	if ok {
		// 计算旧 Token 自然过期时间
		remainingTTL := int(int64(expTimeFloat) - time.Now().Unix())
		if remainingTTL > 0 {
			// 将旧 Token 写进 Redis，存活时间等于剩余物理寿命
			err = l.svcCtx.Redis.AddTokenBlacklist(redisConn, in.RefreshToken, remainingTTL)
			if err != nil {
				l.Logger.Errorf("将旧 Token 写入黑名单失败: %v", err) //不阻断正常运行
			}
		}
	}
	return &user.RefreshResp{
		StatusCode:   0,
		StatusMsg:    "refresh success",
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
