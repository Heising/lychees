package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"lychees-server/configs"
	"lychees-server/logs"
	"lychees-server/models"
	"lychees-server/utils"
	"strconv"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/redis/go-redis/v9"
)

var redisDB *redis.Client

// 哈希表名
var refreshAuth string = "refresh_auth"
var checkBookmarksStatus string = "check_bookmarks_status"
var emailVerifyCode string = "email_verify_code"

// 到期时间是30天后
const expire time.Duration = 30 * 24 * time.Hour
// 邮箱验证码清理时间
const emailCodeExpire time.Duration = 10 * time.Minute

var TokenBlacklist *utils.ExpiringSet

// const expire time.Duration = time.Hour
func initRedisDb() {
	redisDB = redis.NewClient(&redis.Options{
		Addr:     configs.Config.Redis.Host + ":" + strconv.Itoa(configs.Config.Redis.Port),
		Password: configs.Config.Redis.Password, // 没有密码，默认值
		DB:       configs.Config.Redis.DBName,   // 默认DB 0
	})

	err := redisDB.Ping(context.Background()).Err()
	if err != nil {
		logs.Logger.Fatal("Redis错误:", "ERR", err)
	} else {
		logs.Logger.Info("RedisDB连接成功!")
	}


	TokenBlacklist = utils.NewExpiringSet(utils.ExpireJwt)

}

func ADDToken(ctx *gin.Context, user *models.ResponseUser) *models.TokenInfo {

	// 从Redis中检索消息并反序列化
	val, err := redisDB.Get(context.Background(),
		refreshAuth+":"+
			fmt.Sprint(user.ID)).Bytes()

	var newTokenInfo models.TokenInfo

	// 签发一个token
	newTokenInfo.Token = utils.Issuer()

	var tokens models.Auth

	now := time.Now()

	newTokenInfo.ExpireUnix = now.Add(expire).Unix()

	switch {
	case err == nil:
		// 解码
		if unmarshalErr := proto.Unmarshal(val, &tokens); unmarshalErr != nil {
			logs.Logger.Fatal(unmarshalErr)
		}
		logs.Logger.Infof("长度是%d", len(tokens.Infos))

		//先清理到期的token
		cleanExpire(&tokens.Infos)
	default:
		if errors.Is(err, redis.Nil) {
			// key不存在 没有登录过的话
			tokens.Email = user.Email
		} else {
			logs.Logger.Error(err)
		}
	}

	// 追加到最后一位
	tokens.Infos = append(tokens.Infos, &models.Info{
		ExpireUnix: now.Add(expire).Unix(),
		Token:      *newTokenInfo.Token,
		LoginTime:  now.Unix(),
		UserAgent:  ctx.GetHeader("User-Agent"),
		ClientIp:   ctx.GetHeader("x-forwarded-for")},
	)

	// 如果大于10，则移除
	if len(tokens.Infos) > 10 {
		logs.Logger.Infof("长度是%d，即将移除第一个元素", len(tokens.Infos))

		tokens.Infos = tokens.Infos[1:]
	}
	// 重新编码
	data, err := proto.Marshal(&tokens)
	if err != nil {
		logs.Logger.Fatal(err)
	}

	// 存储到redis
	err = redisDB.Set(context.Background(),
		refreshAuth+":"+fmt.Sprint(user.ID),
		data, expire).Err()
	if err != nil {
		logs.Logger.Fatal(err)

	}
	return &newTokenInfo
}

// 验证token是否存在,
func VerifyToken(ctx *gin.Context, userId string, refreshToken *string) (newTokenInfo *models.TokenInfo) {
	// 从Redis中检索消息并反序列化
	val, err := redisDB.Get(context.Background(),
		refreshAuth+":"+userId).Bytes()
	var tokens models.Auth
	now := time.Now()
	switch {
	// redis出现错误
	case err != nil:
		if errors.Is(err, redis.Nil) {
			// 用户从未登录
			return nil
		}
		//不知道什么错误，直接失败
		logs.Logger.Error(err)
		return nil

	default:
		// 解码
		err := proto.Unmarshal(val, &tokens)
		if err != nil {
			logs.Logger.Fatal(err)
			return nil
		}

		//先清理到期的token
		isUpdate := cleanExpire(&tokens.Infos)

		for i := range tokens.Infos {
			logs.Logger.Infof("遍历第%d\n", i)
			logs.Logger.Info((tokens.Infos)[i])
			if *refreshToken == tokens.Infos[i].Token {
				//如果快到期
				if tokens.Infos[i].ExpireUnix < now.Add(24*time.Hour).Unix() {
					logs.Logger.Infof("快到期%d", tokens.Infos[i].ExpireUnix)
					// 更新元素
					newTokenInfo = &models.TokenInfo{
						Token:      utils.Issuer(),
						ExpireUnix: now.Add(expire).Unix(),
					}

					tokens.Infos[i] = &models.Info{
						ExpireUnix: now.Add(expire).Unix(),
						//不应该变更第一次登录时间
						Token:     *newTokenInfo.Token,
						UserAgent: ctx.GetHeader("User-Agent"),
						ClientIp:  ctx.GetHeader("x-forwarded-for"),
					}
					isUpdate = true

					break
				} else {
					//没到期，但匹配到，退出循环

					newTokenInfo = &models.TokenInfo{
						Token:      &tokens.Infos[i].Token,
						ExpireUnix: tokens.Infos[i].ExpireUnix,
					}

					break
				}

			}
		}

		if isUpdate {
			// 重新编码
			data, err := proto.Marshal(&tokens)

			if err != nil {
				logs.Logger.Fatal(err)
				return nil
			}
			// 存储到redis
			//err = redisDB.HSet(context.Background(), refreshAuth, userId, data).Err()
			err = redisDB.Set(context.Background(),
				refreshAuth+":"+userId,
				data, expire).Err()
			// 存储失败
			if err != nil {
				logs.Logger.Error(err)
				return nil
			}
		}
		return newTokenInfo
	}

}

// 清除过期的token 如果有更新，则返回true
func cleanExpire(Infos *[]*models.Info) (isUpdate bool) {

	now := time.Now()
	for i := 0; i < len(*Infos); {

		//如果有到期的
		//先清理到期的token
		logs.Logger.Infof("遍历第%d\n", i+1)
		logs.Logger.Info((*Infos)[i])
		if (*Infos)[i].ExpireUnix < now.Unix() {
			isUpdate = true
			//if i >= len(tokens.Infos) {
			//	break
			//}
			logs.Logger.Infof("%d过期了\n", (*Infos)[i].ExpireUnix)
			TokenBlacklist.Add((*Infos)[i].Token)
			*Infos = append((*Infos)[:i], (*Infos)[i+1:]...)

		} else {

			//如果没有 循环下一个
			i++
		}

	}
	return isUpdate
}

// 退出登录用
func CleanToken(userId string, cleanToken *string) {
	// 从Redis中检索消息并反序列化
	val, err := redisDB.Get(context.Background(),
		refreshAuth+":"+userId).Bytes()
	var tokens models.Auth
	switch {
	// redis出现错误
	case err != nil:
		if errors.Is(err, redis.Nil) {
			// 用户从未登录
			return
		}
		//不知道什么错误，直接失败
		logs.Logger.Error(err)
		return

	default:
		// 解码
		err = proto.Unmarshal(val, &tokens)
		if err != nil {
			logs.Logger.Error(err)
			return
		}

		//先清理到期的token
		isUpdate := cleanExpire(&tokens.Infos)
		logs.Logger.Infof("长度是%d\n", len(tokens.Infos))
		for i := 0; i < len(tokens.Infos); i++ {
			if *cleanToken == tokens.Infos[i].Token {
				logs.Logger.Infof("添加到黑名单%s\n", tokens.Infos[i].Token)
				//添加到黑名单
				TokenBlacklist.Add(tokens.Infos[i].Token)
				//直接删除
				logs.Logger.Infof("移除token: %s", tokens.Infos[i].Token)
				tokens.Infos = append((tokens.Infos)[:i], (tokens.Infos)[i+1:]...)
				isUpdate = true
				break
			}
		}

		if isUpdate {
			// 重新编码
			data, err := proto.Marshal(&tokens)
			if err != nil {
				logs.Logger.Error(err)
			}

			// 存储到redis
			err = redisDB.Set(context.Background(),
				refreshAuth+":"+userId,
				data, expire).Err()
			if err != nil {
				logs.Logger.Error(err)
			}
		}
		return

	}
}

// 清除所有token 给修改密码或者修改邮箱时踢下线用
func CleanAllToken(userId string) {
	// 从Redis中检索消息并反序列化
	val, err := redisDB.Get(context.Background(),
		refreshAuth+":"+userId).Bytes()
	var tokens models.Auth

	// redis出现错误
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 用户从未登录
			return
		}
		//不知道什么错误，直接失败
		logs.Logger.Error(err)
		return
	}

	// 解码，需要拿到token放到黑名单
	err = proto.Unmarshal(val, &tokens)
	if err != nil {
		logs.Logger.Error(err)
		return
	}
	// 初始化一个长度为10的字符串切片
	keys := make([]string, 10)
	logs.Logger.Infof("长度是%d\n", len(tokens.Infos))
	for i := 0; i < len(tokens.Infos); i++ {
		logs.Logger.Infof("遍历第%d\n", i)
		//添加到token黑名单
		keys = append(keys, tokens.Infos[i].Token)
	}
	TokenBlacklist.AddSlice(keys)

	//直接移除
	err = redisDB.Del(context.Background(), refreshAuth+":"+userId).Err()
	if err != nil {
		logs.Logger.Info(err)
	}

}

func UpdateBookmarkStatus(userId string, status int64) {
	err := redisDB.Set(context.Background(),
		checkBookmarksStatus+":"+userId,
		status, expire).Err()

	// 存储失败
	if err != nil {
		logs.Logger.Fatal(err)
	}

}

// 判断是否一样，一样返回true，不存在不一致返还false
func CheckBookmarkStatus(userId uint, status string) bool {
	result, err := redisDB.Get(context.Background(),
		checkBookmarksStatus+":"+fmt.Sprint(userId)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false
		}
		logs.Logger.Fatal(err)
		return false
	}

	if result == status {
		return true
	}

	return false
}

func SetEmailCode(email, code string) (err error) {

	// 存储到redis
	err = redisDB.Set(context.Background(),
		emailVerifyCode+":"+email,
		code, emailCodeExpire).Err()
	if err != nil {
		logs.Logger.Fatal(err)
		return err
	}
	return nil
}

func VerifyEmailCode(email, code string) bool {

	result, err := redisDB.Get(context.Background(), emailVerifyCode+":"+email).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			//找不到邮箱
			return false
		}
		logs.Logger.Error(err)
		return false
	}

	if result == code {
		_, err = redisDB.Del(context.Background(), emailVerifyCode+":"+email).Result()
		logs.Logger.Infof("删除结果是%s", result)
		if err != nil {
			logs.Logger.Error(err)
			return false
		}
		return true
	} else {
		return false
	}

}

func GetDevices(userId string) (tokens *models.Auth) {
	// 从Redis中检索消息并反序列化
	val, err := redisDB.Get(context.Background(),
		refreshAuth+":"+userId).Bytes()
	//初始化，防止出bug
	tokens = &models.Auth{}
	switch {
	// redis出现错误
	case err != nil:
		if errors.Is(err, redis.Nil) {
			// 用户从未登录
			return nil
		}
		//不知道什么错误，直接失败
		logs.Logger.Error(err)
		return nil

	default:
		// 解码
		err := proto.Unmarshal(val, tokens)
		if err != nil {
			logs.Logger.Error(err)
			return nil
		}

		//先清理到期的token
		if cleanExpire(&tokens.Infos) {
			// 重新编码
			data, err := proto.Marshal(tokens)

			if err != nil {
				logs.Logger.Fatal(err)
				return nil
			}
			// 存储到redis
			err = redisDB.Set(context.Background(),
				refreshAuth+":"+userId,
				data, expire).Err()
			// 存储失败
			if err != nil {
				logs.Logger.Fatal(err)
			}

		}

		return tokens

	}
}
