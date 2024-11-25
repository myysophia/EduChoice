package user

import (
	"context"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"strconv"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Province     string `json:"province"`
	ExamType     string `json:"exam_type"`
	SchoolType   string `json:"school_type"`
	Physics      bool   `json:"physics"`
	History      bool   `json:"history"`
	Chemistry    bool   `json:"chemistry"`
	Biology      bool   `json:"biology"`
	Geography    bool   `json:"geography"`
	Politics     bool   `json:"politics"`
	Score        int    `json:"score"`
	ProvinceRank int    `json:"province_rank"`
	Holland      string `json:"holland"`
	Interests    string `json:"interests"`
}

func InsertOne(user *User) error {
	return common.DB.Create(user).Error
}

func UpdateOne(id int, user *User) error {
	tx := common.DB.Begin()
	if err := common.DB.Where("id = ?", id).Updates(user).Error; err != nil {
		tx.Rollback()
		return err
	}
	owner := "/api/v1/zy/mock" + strconv.Itoa(id)
	if err := common.REDIS.Del(context.Background(), owner).Err(); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func FindOne(id int) (*User, error) {
	u := &User{}
	if err := common.DB.First(u, id).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func FindOneEP(email string, password string) (*User, error) {
	u := &User{}
	if err := common.DB.First(u, "email = ? and password = ?", email, password).Error; err != nil {
		return nil, err
	}
	return u, nil
}
