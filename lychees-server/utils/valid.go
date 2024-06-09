package utils

import (
	"net/mail"
	"unicode"
)

// 校验用户名
//func IsUsernameValid(username string) bool {
//	regex := regexp.MustCompile(`^[a-zA-Z0-9_]{4,16}$`)
//	return regex.MatchString(username)
//}

// 校验昵称
func IsNicknameValid(nickname string) bool {

	return len([]rune(nickname)) > 16
}

// 校验密码
func IsPasswordValid(password string) bool {
	//改成只校验长度，不管复杂性 如果用户不使用合法登录客户端，弱密码注册，那直接不管
	if len(password) < 6 || len(password) > 72 {
		return false
	}

	//var (
	//	upper   bool
	//	lower   bool
	//	digit   bool
	//	special bool
	//)
	//for _, c := range password {
	//	// Check special characters
	//	if unicode.IsSpace(c) || unicode.IsControl(c) {
	//		return false
	//	}
	//	if c > 127 {
	//		return false
	//	}
	//	switch {
	//	case c >= 'A' && c <= 'Z':
	//		upper = true
	//	case c >= 'a' && c <= 'z':
	//		lower = true
	//	case c >= '0' && c <= '9':
	//		digit = true
	//	case c == '!' || c == '@' || c == '#' || c == '$' || c == '%' || c == '^' || c == '&' ||
	//		c == '*' || c == '(' || c == ')' || c == ',' || c == '.' || c == '?' || c == '"' ||
	//		c == ';' || c == ':' || c == '{' || c == '}' || c == '|' || c == '<' || c == '>':
	//		special = true
	//	}
	//}
	return true
}

// 校验邮箱格式
func CheckEmail(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil

}

// 检测是不是字符串数字
func CheckDigit(s string) (isInt bool) {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
