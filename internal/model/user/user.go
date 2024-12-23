package user

import (
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"strings"
)

// User 用户模型
type User struct {
	ID           int     `json:"id" gorm:"primaryKey;column:id"`
	Username     string  `json:"username" gorm:"uniqueIndex;not null;column:username"`
	Email        string  `json:"email" gorm:"uniqueIndex;not null;column:email"`
	Password     string  `json:"-" gorm:"not null;column:password"`  // json:"-" 避免返回密码
	Province     *string `json:"province" gorm:"column:province"`
	ExamType     *string `json:"exam_type" gorm:"column:exam_type"`
	SchoolType   *string `json:"school_type" gorm:"column:school_type"`
	Score        *int    `json:"score" gorm:"column:score"`
	ProvinceRank *int    `json:"province_rank" gorm:"column:province_rank"`
	Physics      *bool   `json:"physics" gorm:"column:physics"`
	History      *bool   `json:"history" gorm:"column:history"`
	Chemistry    *bool   `json:"chemistry" gorm:"column:chemistry"`
	Biology      *bool   `json:"biology" gorm:"column:biology"`
	Geography    *bool   `json:"geography" gorm:"column:geography"`
	Politics     *bool   `json:"politics" gorm:"column:politics"`
	Holland      *string `json:"holland" gorm:"column:holland"`
	Interests    *string `json:"interests" gorm:"column:interests"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// GetInterests 获取兴趣爱好数组
func (u *User) GetInterests() []string {
	if u.Interests == nil {
		return []string{}
	}
	return strings.Split(*u.Interests, ",")
}

// IsComplete 检查用户信息是否完整
func (u *User) IsComplete() bool {
	if u.Score == nil || u.ExamType == nil || u.Province == nil || 
	   u.Interests == nil || u.Holland == nil {
		return false
	}
	interests := u.GetInterests()
	return len(interests) > 0
}

// FindOne 根据ID查找用户
func FindOne(id int) (*User, error) {
	var u User
	err := common.DB.Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByEmail 根据邮箱查找用户
func FindByEmail(email string) (*User, error) {
	var u User
	err := common.DB.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByUsername 根据用户名查找用户
func FindByUsername(username string) (*User, error) {
	var u User
	err := common.DB.Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByAccountAndPassword 通过账号(邮箱或用户名)和密码查找用户
func FindByAccountAndPassword(account, password string) (*User, error) {
	var u User
	err := common.DB.Where("(email = ? OR username = ?) AND password = ?",
		account, account, password).First(&u).Error
	return &u, err
}

// UpdateOne 更新用户信息
func UpdateOne(id int, u *User) error {
	return common.DB.Model(&User{ID: id}).Updates(u).Error
}

// InsertOne 插入新用户
func InsertOne(u *User) error {
	return common.DB.Create(u).Error
}
