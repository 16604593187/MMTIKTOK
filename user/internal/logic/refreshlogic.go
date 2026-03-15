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
	// Parse函数校验签名与过期时间
	token, err := jwt.Parse(in.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				l.Logger.Infof("风控拦截：refreshToken 已自然过期")
				// 🌟 重构：不再 return err，而是返回业务错误码 10001
				return &user.RefreshResp{StatusCode: 10001, StatusMsg: "长Token已过期，请重新登录"}, nil
			}
		}
		l.Logger.Errorf("风控拦截：长 Token 签名非法或被篡改: %v", err)
		return &user.RefreshResp{StatusCode: 10002, StatusMsg: "无效的Token"}, nil
	}
	
	if !token.Valid {
		return &user.RefreshResp{StatusCode: 10002, StatusMsg: "无效的Token"}, nil
	}

	// 提取UserId
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &user.RefreshResp{StatusCode: 10003, StatusMsg: "Token 格式解析失败"}, nil
	}

	userIdFloat, ok := claims["userId"].(float64)
	if !ok {
		return &user.RefreshResp{StatusCode: 10003, StatusMsg: "Token 缺少用户信息"}, nil
	}
	userId := uint64(userIdFloat)

	oldJti, ok := claims["jti"].(string)
	if !ok {
		return &user.RefreshResp{StatusCode: 10003, StatusMsg: "非法 Token：缺少设备标识"}, nil
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
		return &user.RefreshResp{StatusCode: 10004, StatusMsg: "账号已被冻结，强制下线"}, nil
	}

	// 生成全新的 Token 组合
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	refreshExpire := l.svcCtx.Config.JwtAuth.RefreshExpire
	newAccessToken, newRefreshToken, newJti, err := getJwtTokens(secretKey, now, accessExpire, refreshExpire, userId)
	if err != nil {
		//系统错误
		return nil, err
	}

	// Lua 轮转脚本 (防重放、防挤占、删旧加新)
	status, err := l.svcCtx.Redis.CheckAndRotateToken(redisConn, userId, oldJti, newJti, remainingTTL)
	if err != nil {
		l.Logger.Errorf("执行 Token 轮转脚本失败: %v", err)
		//Redis 挂了，也是系统错误
		return nil, errors.New("系统繁忙，刷新失败")
	}

	switch status {
	case -1:
		//触发保护性强制下线
		tokenTtl := int(refreshExpire)
		errBan := l.svcCtx.Redis.ForceLogoutAll(redisConn, userId, tokenTtl)
		if errBan == nil {
			l.Logger.Errorf("🚨 保护性下线触发：检测到重放攻击 (jti:%s)，已将用户 %d 的所有设备踢下线", oldJti, userId)
		}
		// 返回业务错误码 10005
		return &user.RefreshResp{
			StatusCode: 10005, 
			StatusMsg: "检测到账号环境异常，为保护您的安全，已退出登录，请重新输入密码",
		}, nil

	case -2:
		//设备拦截
		l.Logger.Errorf("风控拦截：用户 %d 的设备(jti:%s) 已被新设备挤下线", userId, oldJti)
		// 返回业务错误码 10006
		return &user.RefreshResp{
			StatusCode: 10006, 
			StatusMsg: "登录状态已失效（可能由于在其他设备登录），请重新输入密码登录",
		}, nil
	}

	// 全部通过，正常返回
	return &user.RefreshResp{
		StatusCode:   0,
		StatusMsg:    "refresh success",
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}