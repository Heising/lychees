package service

import (
	"errors"
	"lychees-server/dao"
	"lychees-server/logs"
	"lychees-server/models"
	"lychees-server/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	cmap "github.com/orcaman/concurrent-map/v2"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

// json request body.

type userInfo struct {
	//在验证码数据库的id
	CaptchaId string `json:"captchaId" binding:"required"`
	//验证码的值
	VerifyValue string `json:"verifyValue" binding:"required"`
	Email       string `json:"email" binding:"required"`
}
type userRegistrationInfo struct {
	userInfo
	Nickname string `json:"nickName" binding:"required"`
}
type userEmailUpdateInfo struct {
	userInfo
	NewEmail string `json:"newEmail" binding:"required"`
}

// 每秒突发求请求限制
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}
type visitorWithRate struct {
	//使用非官方的并发cmap库
	visitor cmap.ConcurrentMap[string, *visitor]
	rate    rate.Limit
}

// 内存存储 10分钟过期
var captchaStore = base64Captcha.DefaultMemStore

// create base64 encoding captcha
var driver *base64Captcha.DriverDigit

// ip限流器的速率
var ipRt = rate.Every(time.Second)

// 邮箱发送验证码速率
var emailRt = rate.Every(5 * time.Minute)

// 邮箱校验验证码速率
var verifyEmailCodeRt = rate.Every(5 * time.Second)

// IP限制
var ipVisitors *visitorWithRate

// 邮箱申请验证码限制
var emailVisitors *visitorWithRate

// 邮箱验证码校验速率限制，短时间防止枚举破解
var EmailCodeverifyVisitors *visitorWithRate

var ErrRateIp = errors.New("验证码请求频繁")
var ErrRateEmail = errors.New("同一邮箱只能五分钟申请一次邮箱验证码")
var ErrRateEmailVerify = errors.New("同一邮箱只能五秒钟校验一次邮箱验证码")

// 过期时间
const rateExpire = 30 * time.Minute

// 新建一个速率器
func newVisitorWithRate(rateLimit rate.Limit) (visitors *visitorWithRate) {
	visitors = &visitorWithRate{
		visitor: cmap.New[*visitor](),
		rate:    rateLimit,
	}
	return visitors
}
func (v *visitorWithRate) getVisitor(key string) *rate.Limiter {
	i, exists := v.visitor.Get(key)
	if !exists {
		vLimiter := rate.NewLimiter(v.rate, 1)
		v.visitor.Set(key, &visitor{
			limiter:  vLimiter,
			lastSeen: time.Now(),
		})
		return vLimiter
	}

	i.lastSeen = time.Now()
	return i.limiter
}

// 检测可以访问
func (v *visitorWithRate) doCheckLimiter(ip string) bool {
	ipLimiter := v.getVisitor(ip)
	return ipLimiter.Allow()
}

// 统一30分钟没有访问就移除 防止占用内存
func (v *visitorWithRate) cleanup() {
	v.visitor.IterCb(func(key string, value *visitor) {
		if time.Since(value.lastSeen) > rateExpire {
			v.visitor.Remove(key)
		}
	})
}

// 生成实例
func init() {
	// 初始化验证码驱动器为数字字符串
	driver = &base64Captcha.DriverDigit{
		Length:   6,
		Width:    240,
		Height:   80,
		MaxSkew:  1,
		DotCount: 100,
	}
	//初始化限制器
	ipVisitors = newVisitorWithRate(ipRt)
	emailVisitors = newVisitorWithRate(emailRt)
	EmailCodeverifyVisitors = newVisitorWithRate(verifyEmailCodeRt)

	go func() {
		//统一清理
		for range time.Tick(30 * time.Second) {
			ipVisitors.cleanup()
			emailVisitors.cleanup()
			EmailCodeverifyVisitors.cleanup()
		}
	}()
}

// 获取验证码
func GenerateCaptchaHandler(ctx *gin.Context) {
	//检查ip时间
	if !ipVisitors.doCheckLimiter(ctx.GetHeader("x-forwarded-for")) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrRateIp.Error()})
		return
	}

	var captcha = base64Captcha.NewCaptcha(driver, captchaStore)
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		//utils.ServerError(ctx)
		logs.Logger.Info(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	logs.Logger.Info(answer)
	ctx.JSON(http.StatusOK, gin.H{"data": map[string]string{"baseCaptcha": b64s, "captchaId": id}})
}

// 验证是否有效 然后发送邮箱注册验证码
func CaptchaVerifyRegister(ctx *gin.Context) {
	//parse request json body
	var param userRegistrationInfo

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		utils.Warning(ctx)
		logs.Logger.Info("字段错误", ctx.GetHeader("x-forwarded-for"))
		return
	}

	if !utils.CheckEmail(param.Email) {
		//utils.Warning(ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "邮箱非法"})
		return
	}

	//verify the captcha 只要在就删除 不管有没有匹配到
	if captchaStore.Verify(param.CaptchaId, param.VerifyValue, true) {
		//申请公钥
		rsaPair, nanoid := utils.GetPublicKey()
		//检查邮箱限制
		if !emailVisitors.doCheckLimiter(param.Email) {
			//但还是，给他key用来加密
			ctx.JSON(http.StatusBadRequest, gin.H{
				"data": map[string]any{
					"publicKey": rsaPair.PublicKey,
					"nanoid":    nanoid,
				},
				"error": ErrRateEmail.Error(),
			})
			return
		}

		if _, err := dao.FindEmail(&param.Email); err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "邮箱已经注册，请登录！",
			})
			return
		}

		//验证成功，给他key用来加密
		ctx.JSON(http.StatusOK, gin.H{
			"data": map[string]any{
				"publicKey": rsaPair.PublicKey,
				"nanoid":    nanoid,
			},
		})
		//生成数字验证码
		digit := utils.GenerateDigit()
		//设置到redis
		err := dao.SetEmailCode(param.Email, digit)
		if err != nil {
			logs.Logger.Info(err)
			utils.ServerError(ctx)
			return

		}
		//给邮箱发送验证码 由于未注册，使用客户端注册时提交的昵称
		utils.SendEmail(param.Email, param.Nickname, digit)

		return
	}
	//set json response

	ctx.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误！"})

}

// 重置密码验证
func CaptchaVerifyPasswordReset(ctx *gin.Context) {
	//parse request json body
	var param userInfo

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		utils.Warning(ctx)
		logs.Logger.Info("字段错误", ctx.GetHeader("x-forwarded-for"))
		return
	}

	if !utils.CheckEmail(param.Email) {
		//utils.Warning(ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "邮箱非法"})
		return
	}

	//verify the captcha 只要在就删除 不管有没有匹配到
	if captchaStore.Verify(param.CaptchaId, param.VerifyValue, true) {
		//申请公钥
		rsaPair, nanoid := utils.GetPublicKey()
		//检查邮箱限制
		if !emailVisitors.doCheckLimiter(param.Email) {
			//但还是，给他key用来加密
			ctx.JSON(http.StatusBadRequest, gin.H{
				"data": map[string]any{
					"publicKey": rsaPair.PublicKey,
					"nanoid":    nanoid,
				},
				"error": ErrRateEmail.Error(),
			})
			return
		}
		var existingUser *models.User
		if existingUser, err = dao.FindEmail(&param.Email); err != nil {
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

		//验证成功，给他key用来加密
		ctx.JSON(http.StatusOK, gin.H{
			"data": map[string]any{
				"publicKey": rsaPair.PublicKey,
				"nanoid":    nanoid,
			},
		})
		//生成数字验证码
		digit := utils.GenerateDigit()
		//设置到redis
		err := dao.SetEmailCode(param.Email, digit)
		if err != nil {
			logs.Logger.Info(err)
			utils.ServerError(ctx)
			return

		}
		//由于已经注册，判断是否存在，从数据库里面获取昵称
		result, err := dao.FindDocuments(existingUser.ID)
		if err != nil {
			logs.Logger.Info(err)
			utils.ServerError(ctx)
			return
		}
		//给邮箱发送验证码
		utils.SendEmail(existingUser.Email, result.Nickname, digit)

		return
	}
	//set json response

	ctx.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误！"})

}

// 更新邮箱验证
func CaptchaVerifyEmailUpdate(ctx *gin.Context) {
	//parse request json body
	var param userEmailUpdateInfo

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		utils.Warning(ctx)
		logs.Logger.Info("字段错误", ctx.GetHeader("x-forwarded-for"))
		return
	}

	if !utils.CheckEmail(param.Email) {
		//utils.Warning(ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "邮箱非法"})
		return
	}
	if !utils.CheckEmail(param.NewEmail) {
		//utils.Warning(ctx)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "新邮箱非法"})
		return
	}
	if param.Email == param.NewEmail {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "新旧邮箱相同，请修改"})
		return
	}
	//verify the captcha 只要在就删除 不管有没有匹配到
	if captchaStore.Verify(param.CaptchaId, param.VerifyValue, true) {
		//申请公钥
		rsaPair, nanoid := utils.GetPublicKey()
		//检查邮箱限制
		if !emailVisitors.doCheckLimiter(param.NewEmail) {
			//但还是，给他key用来加密
			ctx.JSON(http.StatusBadRequest, gin.H{
				"data": map[string]any{
					"publicKey": rsaPair.PublicKey,
					"nanoid":    nanoid,
				},
				"error": ErrRateEmail.Error(),
			})

			return
		}
		var existingUser *models.User
		if existingUser, err = dao.FindEmail(&param.Email); err != nil {
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

		if _, err := dao.FindEmail(&param.NewEmail); err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "新邮箱已经存在，请更换！",
			})
			return
		}

		//验证成功，给他key用来加密
		ctx.JSON(http.StatusOK, gin.H{
			"data": map[string]any{
				"publicKey": rsaPair.PublicKey,
				"nanoid":    nanoid,
			},
		})
		//生成数字验证码
		digit := utils.GenerateDigit()
		//设置到redis
		err := dao.SetEmailCode(param.NewEmail, digit)
		if err != nil {
			logs.Logger.Info(err)
			utils.ServerError(ctx)
			return

		}
		//由于已经注册，判断是否存在，从数据库里面获取昵称
		result, err := dao.FindDocuments(existingUser.ID)
		if err != nil {
			logs.Logger.Info(err)
			utils.ServerError(ctx)
			return
		}
		//给邮箱发送验证码
		utils.SendEmail(param.NewEmail, result.Nickname, digit)

		return
	}
	//set json response

	ctx.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误！"})

}