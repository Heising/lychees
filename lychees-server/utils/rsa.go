package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"lychees-server/logs"
	"strings"
	"sync"
	"time"
)

type RsaPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  string
	ClearAt    int64
}

var rsaPairs = make(map[string]*RsaPair)
var rsaPairsLock sync.RWMutex

func init() {
	// 添加密钥
	go func() {
		//每30分钟生成一个rsa
		for {
			rsaPairsLock.Lock()

			privateKey, publicKey := GenRsaKey()
			rsaPairs[GenerateSalt()] = &RsaPair{
				PrivateKey: privateKey,
				PublicKey:  publicKey,
				ClearAt:    time.Now().Add(time.Hour).Unix(),
			}
			//由于睡眠30分钟需要解锁
			rsaPairsLock.Unlock()

			time.Sleep(30 * time.Minute)
			//重新加锁
			rsaPairsLock.Lock()
			logs.Logger.Infof("长度是：%d", len(rsaPairs))
			for index := range rsaPairs {
				if rsaPairs[index].ClearAt < time.Now().Unix() {
					delete(rsaPairs, index)
				}

			}

			rsaPairsLock.Unlock()
		}
	}()

}

// 获取最新的键值对信息
func GetPublicKey() (*RsaPair, string) {
	//加只读锁
	rsaPairsLock.RLock()
	defer rsaPairsLock.RUnlock()
	var i string
	var unixStamp = time.Now().Unix()
	//比较最大的
	for index := range rsaPairs {
		if rsaPairs[index].ClearAt > unixStamp {
			unixStamp = rsaPairs[index].ClearAt
			i = index
		}

	}
	return rsaPairs[i], i
}

func GenRsaKey() (privateKey *rsa.PrivateKey, publicKey string) {
	priKey, err2 := rsa.GenerateKey(rand.Reader, 2048)
	if err2 != nil {
		panic(err2)
	}

	puKey := &priKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(puKey)
	if err != nil {
		panic(err)
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubKey := pem.EncodeToMemory(block)
	return priKey, string(pubKey)
}


func RsaDecryptBase64(encryptedData string, nanoid string) (string, error) {
	if v, ok := rsaPairs[nanoid]; ok {
		//base64解码数据
		encryptedDecodeBytes, err := base64.StdEncoding.DecodeString(encryptedData)
		if err != nil {
			return "", err
		}
		//解密数据
		originalBytes, err := rsa.DecryptPKCS1v15(rand.Reader, v.PrivateKey, encryptedDecodeBytes)
		if err != nil {
			return "", err
		}
		originalData := string(originalBytes)
		index := strings.Index(originalData, nanoid)
		var password string
		if index == 0 {
			password = originalData[len(nanoid):]
		} else {
			return "", errors.New("密码不符合规则")
		}

		return password, err
	}
	
	return "", errors.New("找不到密钥")
}
