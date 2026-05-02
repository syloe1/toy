package scheduler

import (
	"log"
	"time"

	"cron/internal/mailer"
	"cron/internal/model"

	// 调度器：负责几点执行、多久执行一次任务
	cronpkg "github.com/robfig/cron/v3"
)

type PendingTaskLister interface {
	//传入状态，返回任务数组和错误。
	ListTasks(status string) ([]model.Task, error)
}

type ReminderJob struct {
	// ✅ 用上面的小接口
	taskService PendingTaskLister
	mailer      mailer.ReminderMailer
}

func NewReminderJob(taskService PendingTaskLister, reminderMailer mailer.ReminderMailer) *ReminderJob {
	return &ReminderJob{
		taskService: taskService,
		mailer:      reminderMailer,
	}
}

func (j *ReminderJob) Run() {
	log.Println("running pending task reminder job")

	tasks, err := j.taskService.ListTasks(model.TaskStatusPending)
	if err != nil {
		log.Printf("list pending tasks failed: %v", err)
		return
	}

	if len(tasks) == 0 {
		log.Println("no pending tasks found, skip sending reminder email")
		return
	}

	if err := j.mailer.SendPendingTasksReminder(tasks); err != nil {
		log.Printf("send reminder email failed: %v", err)
		return
	}

	log.Printf("reminder email sent successfully, pending tasks: %d", len(tasks))
}

// 创建一个 cron 定时器
// spec：定时规则（例如 0 0 9 * * * 每天9点）
// timezone：时区（Asia/Shanghai）
// job：要执行的任务（就是你那个 ReminderJob）

func NewCronScheduler(spec, timezone string, job cronpkg.Job) (*cronpkg.Cron, error) {
	// 1. 加载时区（确保按中国时间跑）
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}
	// 2. 创建一个 cron 实例
	scheduler := cronpkg.New(cronpkg.WithLocation(location))

	// 3. ✅ 核心：把「定时规则」和「任务」绑定！
	// 意思：到 spec 时间，就自动调用 job.Run()
	if _, err := scheduler.AddJob(spec, job); err != nil {
		return nil, err
	}

	return scheduler, nil
}
