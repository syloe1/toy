package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	appconfig "cron/internal/config"
	"cron/internal/dao"
	"cron/internal/handler"
	"cron/internal/mailer"
	"cron/internal/model"
	"cron/internal/router"
	"cron/internal/scheduler"
	"cron/internal/service"
)

func main() {
	configPath := strings.TrimSpace(os.Getenv("CONFIG_PATH"))
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := appconfig.Load(configPath)
	if err != nil {
		log.Fatalf("load config %s: %v", configPath, err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("validate config: %v", err)
	}
	//连接数据库
	db, err := dao.NewMySQLDB(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	if err := db.AutoMigrate(&model.Task{}); err != nil {
		log.Fatalf("auto migrate task table: %v", err)
	}
	//初始化 DAO → Service → Handler
	taskDAO := dao.NewTaskDAO(db)
	taskService := service.NewTaskService(taskDAO)

	//初始化邮件 + 定时任务
	smtpMailer := mailer.NewSMTPMailer(cfg.SMTP, cfg.Reminder.Recipients)
	reminderJob := scheduler.NewReminderJob(taskService, smtpMailer)

	jobScheduler, err := scheduler.NewCronScheduler(cfg.Reminder.CronSpec, cfg.Reminder.Timezone, reminderJob)
	if err != nil {
		log.Fatalf("init cron scheduler: %v", err)
	}
	jobScheduler.Start()
	defer jobScheduler.Stop()

	taskHandler := handler.NewTaskHandler(taskService)

	//初始化路由 + 启动 HTTP 服务
	engine := router.New(taskHandler)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("todo server listening on %s", addr)

	if err := engine.Run(addr); err != nil {
		log.Fatalf("run http server: %v", err)
	}
}
