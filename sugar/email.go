package sugar

import (
	"context"
	"crypto/tls"
	"github.com/qingfeng-lian/gosugar/logger"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

// ================================================================================
// 邮件发送
// ================================================================================

type Email struct {
	ctx  context.Context
	to   []string
	from SendEmailFrom
}

type SendEmailFrom struct {
	From     string `json:"from"`
	Pwd      string `json:"pwd"`
	SmtpHost string `json:"smtp_host"`
	SmtpPort int    `json:"smtp_port"`
}

func NewEmail(ctx context.Context, from SendEmailFrom) Email {
	return Email{
		ctx:  ctx,
		from: from,
	}
}

func (that Email) SetTo(toEmail []string) Email {
	that.to = toEmail
	return that
}

func (that Email) Send(title string, contentType string, body string) error {
	// 创建一个新的消息
	m := gomail.NewMessage()

	// 设置发件人和收件人
	m.SetHeader("From", that.from.From)
	m.SetHeader("To", that.to...)

	// 设置邮件主题
	m.SetHeader("Subject", title)

	// 设置邮件正文为HTML格式
	m.SetBody(contentType, body)

	// 可选：添加附件
	// m.Attach("/path/to/file.txt")

	// 创建拨号器并设置SMTP服务器地址和端口
	d := gomail.NewDialer(that.from.SmtpHost, that.from.SmtpPort, that.from.From, that.from.Pwd)

	// 如果需要使用SSL/TLS加密连接，确保设置正确的端口（如465或587）。
	// 对于163邮箱，你可能需要明确启用TLS:
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // 注意：仅用于测试目的，生产环境中应正确配置TLS。

	// 拨号并发送邮件
	err := d.DialAndSend(m)
	logger.Info(that.ctx, "邮件发送", zap.Any("from", that.from.From), zap.Any("to", that.to), zap.Error(err))
	return err
}
