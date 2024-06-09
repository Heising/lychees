package utils

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"lychees-server/logs"
	"mime"
	"net"
	"net/smtp"
)

var header = make(map[string]string)

// SMTP 服务地址
var host = "***.com"

// TODO:
// SMTP端口号 465（SSL加密）
var port = 465
var email = "noreply@***************.**"
var password = "********************************"

func init() {
	chineseName := "荔枝书签页"
	encodedName := mime.QEncoding.Encode("utf-8", chineseName)
	fromHeader := fmt.Sprintf("%s <%s>", encodedName, email)
	header["From"] = fromHeader
	header["Subject"] = mime.QEncoding.Encode("utf-8", "需要验证您的邮箱地址")
	header["Content-Type"] = "text/html; charset=UTF-8"
}
func SendEmail(toEmail, nickname, code string) {
	header["To"] = toEmail

	var buf bytes.Buffer
	err := verifyEmail(nickname, code).Render(context.Background(), &buf)
	if err != nil {
		logs.Logger.Error(err)
	}
	body := buf.String()

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth(
		"",
		email,
		password,
		host,
	)

	err = SendMailUsingTLS(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		email,
		[]string{toEmail},
		[]byte(message),
	)

	if err != nil {
		logs.Logger.Error(err)
	} else {
		logs.Logger.Info("Send mail success!")
	}

}

// return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {

	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		logs.Logger.Error("Create smpt client error:", err)
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				logs.Logger.Error("Error during AUTH", err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
