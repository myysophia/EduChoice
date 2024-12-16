package router

//func Load(e *gin.Engine) {
//    // API v1 路由组
//    v1 := e.Group("/api/v1")
//    {
//        // 认证相关路由
//        authGroup := v1.Group("/auth")
//        {
//            authGroup.POST("/login", auth.Login)           // 登录
//            authGroup.POST("/register", auth.Register)     // 注册
//            authGroup.GET("/email_code", auth.SendEmailVerificationCode) // 发送验证码
//        }
//
//        // 需要认证的路由组
//        authed := v1.Group("")
//        authed.Use(middleware.JWT())
//        {
//            // 用户相关路由
//            userGroup := authed.Group("/user")
//            {
//                userGroup.POST("/info", user.SetInfo)      // 设置用户信息
//                userGroup.GET("/info", user.GetInfo)       // 获取用户信息
//            }
//
//            // 专业推荐相关路由
//            majorGroup := authed.Group("/major")
//            {
//                majorGroup.GET("/recommend", user.GetRecommend) // 获取推荐
//            }
//        }
//    }
//}
