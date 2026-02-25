package logic

import (
	"context"
	"mini-tiktok/user/internal/logic/utils"
	"mini-tiktok/user/internal/svc"
	"mini-tiktok/user/model"
	"mini-tiktok/user/user"
	"time"

	"github.com/golang-jwt/jwt/v4"
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
	accessToken, refreshToken, err := getJwtTokens(secretKey, now, accessExpire, refreshExpire, userInfo.Id)
	if err != nil {
		return nil, err // 签发失败，抛出系统异常
	}
	return &user.LoginResp{
		StatusCode:   STATUS_SUCCESS,
		StatusMsg:    STATUS_SUCCESS_MSG,
		UserID:       userInfo.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// JWT签发算法
func getJwtTokens(secretKey string, iat, accessExpire, refreshExpire int64, userId uint64) (string, string, error) {
	//生成accessToken
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + accessExpire
	claims["iat"] = iat
	claims["userId"] = userId
	accessTokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := accessTokenClaim.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}
	refreshClaims := make(jwt.MapClaims)
	refreshClaims["exp"] = iat + refreshExpire
	refreshClaims["iat"] = iat
	refreshClaims["userId"] = userId
	refreshTokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshTokenClaim.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
