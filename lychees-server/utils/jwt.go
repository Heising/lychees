package utils

import (
	"errors"
	"lychees-server/configs"
	"lychees-server/logs"
	"lychees-server/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 加密校验的jwt密钥
var mySigningKey = []byte(configs.Config.JwtSigningKey)

//过期时间
const ExpireJwt = 24 * time.Hour

type MyCustomClaims struct {
	ID uint `json:"id"`
	TokenInfo *models.TokenInfo `json:"tokenInfo"`
	jwt.RegisteredClaims
}

func GenerateJWT(user *models.ResponseUser, tokenInfo *models.TokenInfo) *string {

	
	claims := MyCustomClaims{
		user.ID,
		tokenInfo,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ExpireJwt)),
			//签发时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
			//签发人
			Issuer:  "Lychees",
			Subject: "Signing",
		},
	}

	// Create claims while leaving out some of the optional fields
	//claims = MyCustomClaims{
	//	"bar",
	//	jwt.RegisteredClaims{
	//		// Also fixed dates can be used for the NumericDate
	//		ExpiresAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
	//		Issuer:    "test",
	//	},
	//}

	newJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, _ := newJwt.SignedString(mySigningKey)

	return &str
}

// 解析 JWT
func ParseJWT(tokenString string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return mySigningKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			logs.Logger.Info("过期了")

			//过期依然断言token，同时抛出错误
			if claims, ok := token.Claims.(*MyCustomClaims); ok {
				logs.Logger.Info(claims)

				return claims, err
			}
		}

		return nil, err
	}

	//断言token
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrInvalidType

}

// GetUserID 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(c *gin.Context) (uint, error) {

	if claims, exists := c.Get("claims"); !exists {
		logs.Logger.Fatal("系统拿不到claims")
		return 0, errors.New("系统拿不到claims")
	} else {
		logs.Logger.Debug("系统拿到claims")

		waitUse := claims.(*MyCustomClaims)
		return waitUse.ID, nil
	}

}

func GetUserClaims(c *gin.Context) (*MyCustomClaims, error) {

	if claims, exists := c.Get("claims"); !exists {
		logs.Logger.Fatal("系统拿不到claims")
		return nil, errors.New("系统拿不到claims")
	} else {
		logs.Logger.Debug("系统拿到claims")

		waitUse := claims.(*MyCustomClaims)
		return waitUse, nil
	}

}
