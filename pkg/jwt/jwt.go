package jwt

import (
	"fmt"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"github.com/golang-jwt/jwt/v4"
)

func SignToken(uid int) (string, error) {
	secretKey := common.CONFIG.String("jwt.secret_key")
	iat := common.CONFIG.Int("jwt.issuer")
	seconds := common.CONFIG.Int("jwt.expire_seconds")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": iat + seconds,
		"iat": iat,
		"uid": uid,
	})

	return token.SignedString([]byte(secretKey))
}

func Parse(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(common.CONFIG.String("jwt.secret_key")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("parse token: invalid claim type")
	}

	if err = claims.Valid(); err != nil {
		return nil, fmt.Errorf("invalid claims: %v", err)
	}

	return claims, nil
}

func ParseGetUID(tokenString string) (int, error) {
	claims, err := Parse(tokenString)
	if err != nil {
		return 0, nil
	}
	uid, ok := claims["uid"].(float64)
	if !ok {
		return 0, fmt.Errorf("ParseGetUID: uid %v type: %T 类型断言失败", claims["uid"], claims["uid"])
	}
	return int(uid), nil
}
