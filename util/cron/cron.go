package cron

import "github.com/robfig/cron/v3"

var (
	CronClient *cron.Cron
)

func Init() {
	CronClient = cron.New()
}
