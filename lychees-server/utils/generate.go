package utils

import (
	"github.com/jaevor/go-nanoid"
	"lychees-server/logs"
)

// 颁发21位的nano id指针
func Issuer() *string {
	generator, err := nanoid.Standard(21)
	if err != nil {
		logs.Logger.Fatal(err)
	}
	s := generator()
	return &s
}

// 颁发16位的Salt
func GenerateSalt() string {
	generator, err := nanoid.Standard(16)
	if err != nil {
		logs.Logger.Fatal(err)
	}
	return generator()
}

// 生成6位的数字
func GenerateDigit() string {
	generator, err := nanoid.CustomASCII("0123456789", 6)
	if err != nil {
		logs.Logger.Fatal(err)
	}
	return generator()
}
