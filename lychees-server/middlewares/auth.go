package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"lychees-server/dao"
	"lychees-server/logs"
	"lychees-server/utils"
	"net/http"
	"time"
)

// 判断JWT
func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("x-jwt")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "兄弟你jwt令牌呢？？？",
			})
			ctx.Abort()
			return
		}
		logs.Logger.Info("进入鉴权")
		claims, err := utils.ParseJWT(token)

		if err != nil {

			if errors.Is(err, jwt.ErrTokenExpired) {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": "兄弟你jwt令牌过期，请刷新",
				})
				ctx.Abort()
				return
			}
			logs.Logger.Info(err)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			//logs.Logger.Fatal(err)
			ctx.Abort()
			return
		}
		if contains := dao.TokenBlacklist.Contains(*claims.TokenInfo.Token); contains {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "兄弟你的令牌无效，请登录",
			})
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()

	}
}

// 判断NanoID jwt过期 token没过期则进入检测颁发
func RefreshAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("x-jwt")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "兄弟你jwt令牌呢？？？",
			})
			ctx.Abort()
			return
		}

		claims, err := utils.ParseJWT(token)

		if err != nil {
			//判断是jwt过期
			if errors.Is(err, jwt.ErrTokenExpired) {
				//如果到期
				if claims.TokenInfo.ExpireUnix < time.Now().Unix() {
					ctx.JSON(http.StatusUnauthorized, gin.H{
						"error": "token令牌到期，请重新登录",
					})
					ctx.Abort()
					return
				}

				ctx.Set("claims", claims)
				ctx.Next()
				return
			}
			//不知道啥错误，直接打印
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			//logs.Logger.Fatal(err)
			ctx.Abort()
			return
		}

		ctx.Status(http.StatusNoContent)
		ctx.Abort()

	}
}
