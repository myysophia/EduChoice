package auth

import (
	"context"
	"fmt"
	"github.com/big-dust/DreamBridge/internal/model/user"
	"github.com/big-dust/DreamBridge/internal/pkg/common"
	"github.com/big-dust/DreamBridge/pkg/jwt"
	"github.com/big-dust/DreamBridge/pkg/utils"
	mail "github.com/xhit/go-simple-mail/v2"
	"go.uber.org/zap"
	"time"
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

func LoginGetToken(account string, password string) (string, error) {
	// 临时账号验证
	if account == "admin" && password == "admin" {
		return jwt.SignToken(1) // 使用ID 1作为临时管理员ID
	}

	// 尝试��过邮箱或用户名登录
	u, err := user.FindByAccountAndPassword(account, password)
	if err != nil {
		return "", err
	}
	return jwt.SignToken(u.ID)
}

// GetUID 通过邮箱和密码获取用户ID
func GetUID(email string, password string) (int, error) {
	// 使用 FindByAccountAndPassword 替代 FindOneEP
	u, err := user.FindByAccountAndPassword(email, password)
	if err != nil {
		return 0, err
	}
	return u.ID, nil
}

func SendEmailCode(email string) error {
	// 生成验证码
	code := utils.GenerateRandomCode(6)
	
	// 创建邮件客户端
	server := mail.NewSMTPClient()
	server.Host = common.CONFIG.String("email.host")
	server.Port = common.CONFIG.Int("email.port")
	server.Username = common.CONFIG.String("email.username")
	server.Password = common.CONFIG.String("email.password")
	server.Encryption = mail.EncryptionTLS  // 使用TLS加密
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	
	// 连接到SMTP服务器
	smtpClient, err := server.Connect()
	if err != nil {
		common.LOG.Error("SMTP Connect Error", zap.Error(err))
		return err
	}
	
	// 创建邮件
	email_msg := mail.NewMSG()
	email_msg.SetFrom(common.CONFIG.String("email.from"))
	email_msg.AddTo(email)
	email_msg.SetSubject("DreamBridge 验证码")
	
	// 设置邮件内容
	email_msg.SetBody(mail.TextHTML, fmt.Sprintf(`
		<h1>DreamBridge 验证码</h1>
		<p>您的验证码是：<strong>%s</strong></p>
		<p>验证码有效期为5分钟，请尽快使用。</p>
		<p>如果这不是您的操作，请忽略此邮件。</p>
	`, code))
	
	// 发送邮件
	if err := email_msg.Send(smtpClient); err != nil {
		common.LOG.Error("Send Email Error", zap.Error(err))
		return err
	}
	
	// 将验证码存入Redis，设置5分钟过期
	key := fmt.Sprintf("email_code:%s", email)
	err = common.REDIS.Set(context.Background(), key, code, 5*time.Minute).Err()
	if err != nil {
		common.LOG.Error("Save Code To Redis Error", zap.Error(err))
		return err
	}
	
	return nil
}
