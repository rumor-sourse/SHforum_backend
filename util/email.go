package util

import (
	"github.com/jordan-wright/email"
	"net/smtp"
	"strings"
)

const (
	emailName     = "chunkai_Wang@qq.com"
	emailPassword = "cbhrivwvqcxaecgb"
	emailHost     = "smtp.qq.com"

	emailHTML       = "<!DOCTYPE html>\n<html>\n  <head>\n    <meta charset=\"UTF-8\">\n    <title>您的验证码</title>\n    <style>\n      body {\n        font-family: Arial, Helvetica, sans-serif;\n        font-size: 14px;\n        line-height: 1.5;\n        color: #333;\n      }\n      .container {\n        max-width: 600px;\n        margin: 0 auto;\n        padding: 20px;\n        border: 1px solid #ccc;\n      }\n      .header {\n        background-color: #f5f5f5;\n        padding: 10px 20px;\n        border-top-left-radius: 5px;\n        border-top-right-radius: 5px;\n      }\n      .header h1 {\n        margin: 0;\n        font-size: 24px;\n        color: #333;\n      }\n      .body {\n        padding: 20px;\n        background-color: #fff;\n        border-bottom-left-radius: 5px;\n        border-bottom-right-radius: 5px;\n      }\n      .code {\n        margin: 20px 0;\n        font-size: 36px;\n        font-weight: bold;\n        color: #00bfff;\n      }\n    </style>\n  </head>\n  <body>\n    <div class=\"container\">\n      <div class=\"header\">\n        <h1>您的验证码</h1>\n      </div>\n      <div class=\"body\">\n        <p>尊敬的用户，</p>\n        <p>感谢您使用我们的服务。您的验证码为：</p>\n        <p class=\"code\">{{code}}</p>\n        <p>请你在{{ExpireTime}}分钟之内进行注册</p>\n      <p>如果您没有进行此操作，请忽略此邮件。</p>\n        <p>祝您使用愉快！</p>\n      </div>\n    </div>\n  </body>\n</html>"
	EmailExpireTime = "5"
)

func getEmailAuth() smtp.Auth {
	return smtp.PlainAuth("", emailName, emailPassword, emailHost)
}

func SendEmailWithCode(to []string, code string) error {
	e := email.NewEmail()
	e.From = "沙河论坛系统 <" + emailName + ">"
	e.To = to
	e.Subject = "沙河论坛验证码"
	getEmailAuth()
	tempHTML := strings.Replace(emailHTML, "{{code}}", code, -1)
	tempHTML = strings.Replace(tempHTML, "{{ExpireTime}}", EmailExpireTime, -1)
	e.HTML = []byte(tempHTML)
	err := e.Send(emailHost+":25", getEmailAuth())
	if err != nil {
		return err
	}
	return nil
}
