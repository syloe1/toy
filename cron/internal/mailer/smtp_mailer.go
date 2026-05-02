package mailer

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	appconfig "cron/internal/config"
	"cron/internal/model"
)

type ReminderMailer interface {
	//你必须能发送「待办任务提醒邮件」
	SendPendingTasksReminder(tasks []model.Task) error
}

type SMTPMailer struct {
	// 邮箱配置（host/port/username/password）
	cfg        appconfig.SMTPConfig
	recipients []string // 收件人列表
}

// 依赖注入：
func NewSMTPMailer(cfg appconfig.SMTPConfig, recipients []string) *SMTPMailer {
	return &SMTPMailer{
		cfg:        cfg,
		recipients: recipients,
	}
}

func (m *SMTPMailer) SendPendingTasksReminder(tasks []model.Task) error {
	if len(tasks) == 0 {
		return nil
	}
	// 1. 建立SMTP认证
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	// 2. 拼接地址：smtp.qq.com:587
	address := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)
	// 3. 构造完整邮件内容（头+体）
	message := BuildPendingTasksMessage(m.cfg.From, m.recipients, tasks)
	// 4. 调用系统库发邮件
	return smtp.SendMail(address, auth, m.cfg.From, m.recipients, message)
}

func BuildPendingTasksMessage(from string, recipients []string, tasks []model.Task) []byte {
	subject := fmt.Sprintf("TodoList Pending Tasks - %s", time.Now().Format("2006-01-02"))
	body := BuildPendingTasksBody(tasks)

	var builder bytes.Buffer
	builder.WriteString(fmt.Sprintf("From: %s\r\n", from))
	builder.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(recipients, ",")))
	builder.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	builder.WriteString("MIME-Version: 1.0\r\n")
	builder.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	builder.WriteString("\r\n")
	builder.WriteString(body)

	return builder.Bytes()
}

func BuildPendingTasksBody(tasks []model.Task) string {
	var builder strings.Builder

	builder.WriteString("You still have unfinished tasks:\n\n")
	for index, task := range tasks {
		builder.WriteString(fmt.Sprintf("%d. %s\n", index+1, task.Title))
		if task.Content != "" {
			builder.WriteString(fmt.Sprintf("   %s\n", task.Content))
		}
		builder.WriteString(fmt.Sprintf("   Created At: %s\n\n", task.CreatedAt.Format("2006-01-02 15:04:05")))
	}

	builder.WriteString("Please complete them as soon as possible.\n")
	return builder.String()
}
