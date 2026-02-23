package logic

import (
	"context"
	"mini-tiktok/user/internal/logic/utils"
	"mini-tiktok/user/internal/svc"
	"mini-tiktok/user/model"
	"mini-tiktok/user/user"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
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
	return &user.LoginResp{
		StatusCode: STATUS_SUCCESS,
		StatusMsg:  STATUS_SUCCESS_MSG,
		UserID:     userInfo.Id,
	}, nil
}
