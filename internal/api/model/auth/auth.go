package auth

// LoginReq 登录请求
type LoginReq struct {
    Account  string `json:"email" binding:"-"` // 移除所有验证规则
    Password string `json:"password" binding:"-"` // 移除所有验证规则
}

// RegisterReq 注册请求
type RegisterReq struct {
    Username string `json:"username" binding:"required,min=2,max=20"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6,max=20"`
    Code     string `json:"code" binding:"required,len=6"`
} 