package auth

import (
	"context"
	"github.com/big-dust/DreamBridge/internal/model/user"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"github.com/big-dust/DreamBridge/pkg/jwt"
)

func OkEmailCode(email string, code string) bool {
	key := common.Dream + "/" + email
	value, err := common.REDIS.Get(context.Background(), key).Result()
	if err != nil {
		return false
	}
	if value != code {
		return false
	}
	if common.REDIS.Del(context.Background(), key).Err() != nil {
		common.LOG.Error("Del Key Failed:" + key)
	}
	return true
}

func Register(username string, email string, password string) error {
	u := &user.User{
		Username: username,
		Email:    email,
		Password: password,
	}
	return user.InsertOne(u)
}

func LoginGetToken(email string, password string) (string, error) {
	uid, err := GetUID(email, password)
	if err != nil {
		return "", err
	}
	return jwt.SignToken(uid)
}

func GetUID(email string, password string) (int, error) {
	u, err := user.FindOneEP(email, password)
	return u.ID, err
}
