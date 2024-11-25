package router

import "github.com/gin-gonic/gin"

func Load(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	useAuth(v1)
	useUser(v1)
	useZY(v1)
}
