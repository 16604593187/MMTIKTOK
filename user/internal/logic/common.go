package logic

import (
	"github.com/golang-jwt/jwt/v4"
)

const (
	STATUS_SUCCESS            = "0"
	STATUS_SUCCESS_MSG        = "OK"
	STATUS_FAIL               = "1"
	STATUS_FAIL_PARAM_MSG     = "Param incorrect"
	STATUS_USER_EXISTS_MSG    = "Username already exists"
	STATUS_USER_NOTEXIST_MSG  = "User not exist"
	STATUS_WRONG_PASSWORD_MSG = "Wrong Password"
	COUNT_NOT_FOUND           = int64(-1)
	OP_FOLLOW                 = "1"
	OP_CANCEL_FOLLOW          = "2"
)
//Jwt签发算法
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
	//生成refreshToken
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
