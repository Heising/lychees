package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 非法客户端警告
func Warning(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "使用非法客户端，将记录IP地址！"})
}
func ServerError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务存在错误，请联系网站管理员"})
}
