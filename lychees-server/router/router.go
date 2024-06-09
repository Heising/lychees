package router

import (
	"github.com/gin-gonic/gin"
	"lychees-server/configs"
	"lychees-server/logs"
	"lychees-server/middlewares"
	"lychees-server/service"
	"strconv"
)

func Start() {

	engine := gin.Default()
	// 接口耗时监控
	engine.Use(middlewares.HandleEndpointLantency())
	engine.Use(middlewares.Cors())
	//engine.GET("/", service.Hello)
	//注册
	engine.POST("/signup", service.SignUp)
	//登录
	engine.POST("/signin", service.SignIn)
	//重置密码
	engine.POST("/password_reset", service.PasswordReset)
	//更新邮箱
	engine.POST("/email_update", service.EmailUpdate)

	//获取加密公钥
	engine.GET("/getkey", service.GetPublicKey)

	//验证码校验
	engine.GET("/captcha", service.GenerateCaptchaHandler)
	engine.POST("/verify_captcha_register", service.CaptchaVerifyRegister)
	engine.POST("/verify_captcha_password_reset", service.CaptchaVerifyPasswordReset)
	engine.POST("/verify_captcha_email_update", service.CaptchaVerifyEmailUpdate)

	//刷新jwt或者token
	engine.GET("/refresh", middlewares.RefreshAuth(), service.RefreshToken)
	//只刷新jwt
	engine.GET("/refreshjwt", middlewares.JWTAuth(), service.RefreshJwt)

	//需要jwt鉴权组
	auth := engine.Group("/auth", middlewares.JWTAuth())

	{
		auth.GET("/devices", service.GetDevices)
		auth.PUT("/personalinfo", service.UpdatePersonalInfo)

		auth.PUT("/iconfont", service.UpdateIconfont)
		auth.DELETE("/logout", service.LogOut)
		auth.DELETE("/device", service.ExpelDevice)
		auth.DELETE("/devices", service.ExpelDevices)

		//书签组
		bookmarks := auth.Group("/bookmarks")
		{
			//查找标签
			bookmarks.GET("/find/:updateAt", service.FindDocuments)

			//增加一个书签
			bookmarks.POST("/push/:index", service.PushSubDocument)

			//修改一个书签
			bookmarks.PUT("/update/:index/:target", service.UpdateSubDocument)
			bookmarks.PUT("/move/:index/:target", service.MoveSubDocument)

			//交换书签
			bookmarks.PUT("/swap/:index", service.SwapSubDocument)

			//移除一个书签
			bookmarks.DELETE("/:rowIndex/:colIndex", service.RemoveSubDocument)
			//移除页面 但不会判断是否为空
			bookmarks.DELETE("/page", service.ClearEmptyArray)

		}
	}

	//bookmarks := engine.Group("/bookmarks", middlewares.JWTAuth())
	err := engine.Run(":" + strconv.Itoa(configs.Config.Port))
	if err != nil {
		logs.Logger.Fatal(err)
	}

}
