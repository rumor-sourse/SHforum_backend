package cron

import "github.com/robfig/cron/v3"

var (
	cronClient *cron.Cron
)

func Init() {
	cronClient = cron.New()
}
