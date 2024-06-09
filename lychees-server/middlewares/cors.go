package middlewares

import (
	"lychees-server/configs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		//origin := c.Request.Header.Get("Origin")
		c.Header("Server", DEFAULT_NAME)

		// fmt.Println("生效")
		c.Header("Access-Control-Allow-Origin", configs.Config.DomainName) // 允许的域名
		c.Header("Access-Control-Allow-Methods",
			"POST, GET, OPTIONS, PUT, DELETE, UPDATE") // 允许的HTTP方法
		c.Header("Access-Control-Allow-Headers",
			"Origin, X-Requested-With, Content-Type, Accept, Authorization, X-jwt") // 允许的请求标头
		c.Header("Access-Control-Expose-Headers",
			"Content-Length, "+
				"Access-Control-Allow-Origin, "+
				"Access-Control-Allow-Headers, "+
				"Cache-Control, "+
				"Content-Language, "+
				"Content-Type,X-Jwt") // 暴露的响应标头
		c.Header("Access-Control-Allow-Credentials", "true") // 是否允许发送凭证（如cookies）

		// OPTIONS方法
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
