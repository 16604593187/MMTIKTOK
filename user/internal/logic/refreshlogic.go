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

	oldJti, ok := claims["jti"].(string)
	if !ok {
		return nil, errors.New("非法 Token：缺少设备标识")
	}

	expTimeFloat, _ := claims["exp"].(float64)
	remainingTTL := int(int64(expTimeFloat) - time.Now().Unix())
	if remainingTTL < 0 {
		remainingTTL = 0
	}

	redisConn := l.svcCtx.Redis.NewRedisConn()
	defer redisConn.Close()

	// 封禁账户检查
	banned, err := l.svcCtx.Redis.IsUserBanned(redisConn, userId)
	if err == nil && banned {
		return nil, errors.New("账号已被冻结，强制下线")
	}

	// 生成全新的 Token 组合
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	refreshExpire := l.svcCtx.Config.JwtAuth.RefreshExpire
	newAccessToken, newRefreshToken, newJti, err := getJwtTokens(secretKey, now, accessExpire, refreshExpire, userId)
	if err != nil {
		return nil, err
	}

	// Lua 轮转脚本 (防重放、防挤占、删旧加新)
	success, err := l.svcCtx.Redis.CheckAndRotateToken(redisConn, userId, oldJti, newJti, remainingTTL)
	if err != nil {
		l.Logger.Errorf("执行 Token 轮转脚本失败: %v", err)
		return nil, errors.New("系统繁忙，刷新失败")
	}

	// 拦截逻辑：重放攻击（黑名单有记录），或者被新设备挤下线，
	if !success {
		l.Logger.Errorf("风控拦截：用户 %d 的设备(jti:%s) 尝试重放或已被挤下线", userId, oldJti)
		return nil, errors.New("登录状态已失效（可能由于在其他设备登录），请重新输入密码登录")
	}
	return &user.RefreshResp{
		StatusCode:   0,
		StatusMsg:    "refresh success",
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
