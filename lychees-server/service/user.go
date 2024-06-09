package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"lychees-server/dao"
	"lychees-server/logs"
	"lychees-server/models"
	"lychees-server/utils"
	"net/http"
	"strconv"
)

//func Hello(ctx *gin.Context) {
//	ctx.Data(200, "text/plain; charset=utf-8", []byte("Hello lychees"))
//}

// SignUp 用户注册
func SignUp(ctx *gin.Context) {
	var registerUser models.RegisterUser
	if err := ctx.ShouldBindJSON(&registerUser); err != nil {
		utils.Warning(ctx)
		logs.Logger.Info("字段错误", ctx.GetHeader("x-forwarded-for"))
		return
	}
	if !dao.VerifyEmailCode(registerUser.Email, registerUser.VerifyCode) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "邮箱验证码错误"})
		return
	}
	if utils.IsNicknameValid(registerUser.Nickname) {
		//utils.Warning(ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "昵称非法"})
		return
	}
	logs.Logger.Info("昵称检测成功")

	//var user = registerUser.LoginUser

	//if !utils.CheckEmail(registerUser.Email) {
	//	utils.Warning(ctx)
	//	ctx.JSON(409, gin.H{
	//		"error": "邮箱非法"})
	//	return
	//}
	//logs.Logger.Info("邮箱检测成功")

	////TODO: 暂时使用标准base64解码
	//decodedData, err := base64.StdEncoding.DecodeString(user.Encrypted)
	//if err != nil {
	//	ctx.JSON(409, gin.H{
	//		"error": "密码非base64编码"})
	//	return
	//}
	//user.Encrypted = string(decodedData)
	//rsa解密
	password, err := utils.RsaDecryptBase64(registerUser.Encrypted, registerUser.Nanoid)
	if err != nil {
		utils.Warning(ctx)
		logs.Logger.Infof("密码解密失败：%s", err)
		return
	}
	if !utils.IsPasswordValid(password) {
		utils.Warning(ctx)
		logs.Logger.Info("密码非法", ctx.GetHeader("x-forwarded-for"))
		return
	}
	logs.Logger.Info("密码检测成功")

	if _, err := dao.FindEmail(&registerUser.Email); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "邮箱已经注册，请登录！",
		})
		return
	}
	logs.Logger.Info("没有已注册")
	var user = models.User{
		Email: registerUser.Email,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
		return
	}
	logs.Logger.Info("密码已加密")
	//密码加密
	user.Encrypted = string(hashedPassword)

	if err := dao.AddUser(&user); err != nil {
		ctx.JSON(500, gin.H{
			"error": "内部服务器错误"})
		return
	}
	logs.Logger.Info("用户已添加")
	registerUser.ID = user.ID
	result, err := dao.DocumentsInit(&registerUser)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "初始化书签失败，请联系网站管理员",
			"user":  &user.ID,
		})
		return
	}

	ctx.JSON(200,
		gin.H{

			"data": map[string]any{
				"message": "注册成功，请登录",
				//"userId":  &user.ID,
				"userId": &result.InsertedID,
			},
		})

}

// 用户登录
func SignIn(ctx *gin.Context) {
	var loginUser models.LoginUser
	var err error

	if err = ctx.ShouldBindJSON(&loginUser); err != nil {
		utils.Warning(ctx)
		logs.Logger.Info("字段错误", ctx.GetHeader("x-forwarded-for"))
		return
	}

	if !utils.CheckEmail(loginUser.Email) {
		utils.Warning(ctx)
		logs.Logger.Info("邮箱非法", ctx.GetHeader("x-forwarded-for"))
		return
	}

	password, err := utils.RsaDecryptBase64(loginUser.Encrypted, loginUser.Nanoid)
	if err != nil {
		utils.Warning(ctx)
		logs.Logger.Infof("密码解密失败：%s", err)
		return
	}

	if !utils.IsPasswordValid(password) {
		utils.Warning(ctx)
		logs.Logger.Info("密码非法", ctx.GetHeader("x-forwarded-for"))
		return
	}

	var existingUser *models.ResponseUser
	user := models.User{
		Email: loginUser.Email,
	}

	if existingUser, err = dao.FindFilterUser(&user); err != nil {
		//logs.Logger.Info(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
			return
		}
		//未知错误
		utils.ServerError(ctx)
		logs.Logger.Fatal(err)
		return

	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Encrypted), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "邮箱或密码错误"})
		return
	}

	tokenInfo := dao.ADDToken(ctx, existingUser)
	newJwt := utils.GenerateJWT(existingUser, tokenInfo)
	// 登录成功
	ctx.JSON(http.StatusOK, gin.H{
		"data": map[string]any{
			"token":    newJwt,
			"userInfo": existingUser,
		},
	})

}

// 重置密码
func PasswordReset(ctx *gin.Context) {
	//防止枚举破解
	if !EmailCodeverifyVisitors.doCheckLimiter(ctx.GetHeader("x-forwarded-for")) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrRateEmailVerify.Error()})
		return
	}
	var passwordResetUser models.PasswordResetUser
	var err error

	if err = ctx.ShouldBindJSON(&passwordResetUser); err != nil {
		utils.Warning(ctx)
		logs.Logger.Info("字段错误", ctx.GetHeader("x-forwarded-for"))
		return
	}

	if !dao.VerifyEmailCode(passwordResetUser.Email, passwordResetUser.VerifyCode) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "邮箱验证码错误"})
		return
	}

	password, err := utils.RsaDecryptBase64(passwordResetUser.Encrypted, passwordResetUser.Nanoid)
	if err != nil {
		utils.Warning(ctx)
		logs.Logger.Infof("密码解密失败：%s", err)
		return
	}

	if !utils.IsPasswordValid(password) {
		utils.Warning(ctx)
		logs.Logger.Info("密码非法", ctx.GetHeader("x-forwarded-for"))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
		return
	}
	logs.Logger.Info("密码已加密")
	//密码加密
	encrypted := string(hashedPassword)

	var existingUser *models.User

	if existingUser, err = dao.FindEmail(&passwordResetUser.Email); err != nil {
		//logs.Logger.Info(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
			return
		}
		//未知错误
		utils.ServerError(ctx)
		logs.Logger.Error(err)
		return

	}
	err = dao.PasswordReset(existingUser, encrypted)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"data": map[string]any{
				"error": "修改密码失败",
			},
		})
		return
	}
	//清除所有token
	dao.CleanAllToken(strconv.Itoa(int(existingUser.ID)))
	// 登录成功
	ctx.JSON(http.StatusOK, gin.H{
		"data": map[string]any{
			"message": "修改密码成功，请登录",
		},
	})

}

// 更新用户邮箱
func EmailUpdate(ctx *gin.Context) {
	var emailUser models.EmailUpdater
	var err error

	if err = ctx.ShouldBindJSON(&emailUser); err != nil {
		utils.Warning(ctx)
		logs.Logger.Info("字段错误", ctx.GetHeader("x-forwarded-for"))
		return
	}

	if !utils.CheckEmail(emailUser.Email) {
		utils.Warning(ctx)
		logs.Logger.Info("邮箱非法", ctx.GetHeader("x-forwarded-for"))
		return
	}

	if !utils.CheckEmail(emailUser.NewEmail) {
		utils.Warning(ctx)
		logs.Logger.Info("新邮箱非法", ctx.GetHeader("x-forwarded-for"))
		return
	}

	if emailUser.Email == emailUser.NewEmail {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "新旧邮箱相同，请修改"})
		return
	}

	if !dao.VerifyEmailCode(emailUser.NewEmail, emailUser.VerifyCode) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "邮箱验证码错误"})
		return
	}
	password, err := utils.RsaDecryptBase64(emailUser.Encrypted, emailUser.Nanoid)
	if err != nil {
		utils.Warning(ctx)
		logs.Logger.Infof("密码解密失败：%s", err)
		return
	}

	if !utils.IsPasswordValid(password) {
		utils.Warning(ctx)
		logs.Logger.Info("密码非法", ctx.GetHeader("x-forwarded-for"))
		return
	}

	var existingUser *models.ResponseUser
	user := models.User{
		Email: emailUser.Email,
	}

	if existingUser, err = dao.FindFilterUser(&user); err != nil {
		//logs.Logger.Info(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
			return
		}
		//未知错误
		utils.ServerError(ctx)
		logs.Logger.Fatal(err)
		return

	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Encrypted), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "邮箱或密码错误"})
		return
	}

	//tokenInfo := dao.ADDToken(ctx, existingUser)
	//newJwt := utils.GenerateJWT(existingUser, tokenInfo)

	if _, err := dao.FindEmail(&emailUser.NewEmail); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "新邮箱已经存在，请更换！",
		})
		return
	}
	//更新邮箱
	err = dao.EmailUpdate(&emailUser)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "修改邮箱失败",
		})
		return
	}
	//清除所有token
	dao.CleanAllToken(strconv.Itoa(int(existingUser.ID)))
	// 返回更新成功的状态码
	ctx.JSON(http.StatusOK, gin.H{
		"data": map[string]any{
			"message": "修改邮箱成功，请登录",
		},
	})

}
func GetDevices(ctx *gin.Context) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.ServerError(ctx)
		return
	}
	devices := dao.GetDevices(strconv.Itoa(int(userID)))

	if devices == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "登录数据没有查询到",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": map[string]any{
			"devices": devices.Infos,
		},
	})

}

// 自己退出登录
func LogOut(ctx *gin.Context) {
	claims, err := utils.GetUserClaims(ctx)
	if err != nil {
		utils.ServerError(ctx)
		return

	}

	dao.CleanToken(strconv.Itoa(int(claims.ID)), claims.TokenInfo.Token)

	// 退出登录成功
	ctx.Status(http.StatusNoContent)

}

// 踢下线单个设备
func ExpelDevice(ctx *gin.Context) {
	var tokenInfo models.Token
	var err error
	if err = ctx.ShouldBindJSON(&tokenInfo); err != nil {
		utils.Warning(ctx)
		logs.Logger.Info("字段错误", ctx.GetHeader("x-forwarded-for"))
		return
	}
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	dao.CleanToken(strconv.Itoa(int(userID)), &tokenInfo.Token)

	// 退出登录成功
	ctx.Status(http.StatusNoContent)

}

// 全部设备踢下线
func ExpelDevices(ctx *gin.Context) {
	claims, err := utils.GetUserClaims(ctx)
	if err != nil {
		utils.ServerError(ctx)
		return

	}

	dao.CleanAllToken(strconv.Itoa(int(claims.ID)))

	// 退出登录成功
	ctx.Status(http.StatusNoContent)

}

// 获取加密公钥
func GetPublicKey(ctx *gin.Context) {
	rsaPair, nanoid := utils.GetPublicKey()
	ctx.JSON(http.StatusOK, gin.H{
		"data": map[string]any{
			"publicKey": rsaPair.PublicKey,
			"nanoid":    nanoid,
		},
	})
}
