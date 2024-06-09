package service

import (
	"github.com/golang-jwt/jwt/v5"
	"lychees-server/dao"
	"lychees-server/logs"
	"lychees-server/models"
	"lychees-server/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 续约Token
func RefreshToken(ctx *gin.Context) {
	claims, err := utils.GetUserClaims(ctx)
	if err != nil {
		utils.ServerError(ctx)
		return

	}
	//到了缓冲时间 可以刷新
	if claims.ExpiresAt.Unix()-time.Now().Unix() < int64(time.Hour.Seconds()) {
		logs.Logger.Info("到了缓冲时间 可以刷新")

		//先验证token是否有效
		newTokenInfo := dao.VerifyToken(ctx, strconv.Itoa(int(claims.ID)), claims.TokenInfo.Token)

		if newTokenInfo != nil {

			existingUser := &models.ResponseUser{
				ID: claims.ID,
			}
			// jwt 需要签发
			newJwt := utils.GenerateJWT(existingUser, newTokenInfo)
			ctx.Header("x-jwt", *newJwt)

			ctx.Status(http.StatusNoContent)
			return

		} else {
			//当查询不到，当过期处理
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": jwt.ErrTokenExpired.Error(),
			})
			return

		}

	}
	//没到缓冲刷新时间，啥也不干
	ctx.Status(http.StatusNoContent)
}

// 续约JWT
func RefreshJwt(ctx *gin.Context) {

	claims, err := utils.GetUserClaims(ctx)
	if err != nil {
		utils.ServerError(ctx)
		return

	}
	//到了缓冲时间 可以刷新
	if claims.ExpiresAt.Unix()-time.Now().Unix() < int64(2*time.Hour.Seconds()) {
		// jwt 需要重新签发
		newJwt := utils.GenerateJWT(
			&models.ResponseUser{
				ID: claims.ID,
			},
			&models.TokenInfo{
				Token:      claims.TokenInfo.Token,
				ExpireUnix: claims.TokenInfo.ExpireUnix,
			},
		)
		ctx.Header("x-jwt", *newJwt)
	}

	//没到缓冲刷新时间，啥也不干
	ctx.Status(http.StatusNoContent)

}
