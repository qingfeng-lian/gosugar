package sugar

import (
	"context"
	"fmt"
	"testing"
)

func TestSendEmail(t *testing.T) {
	contentType := "text/html"
	body := "<h1>Hello!!</h1><p>aaaaaaaaaaaaa.</p>"
	err := NewEmail(context.Background(), SendEmailFrom{
		From:     "zhize02@163.com",
		Pwd:      "MHTVD3nwUfVKs8vB",
		SmtpHost: "smtp.163.com",
		SmtpPort: 465,
	}).SetTo([]string{"471439353@qq.com"}).Send("我是标题", contentType, body)
	fmt.Println("邮件发送", err)
}
