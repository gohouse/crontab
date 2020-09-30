package crontab

import "context"

type TaskObject struct {
	cron     *CronTab
	callback HandleFunc
	args     []interface{}
	cancel   context.CancelFunc
	title    string
	taskId   string
}

func (so *TaskObject) start() {
	if so.cron.running == true {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	so.cron.ctx = ctx
	so.cancel = cancel
	so.cron.Run(so.callback, so.args...)
}
func (so *TaskObject) stop() {
	if so.cron.running == false {
		return
	}
	(so.cancel)()
}
func (so *TaskObject) IsRunning() bool {
	return so.cron.IsRunning()
}
func (so *TaskObject) Title() string {
	return so.title
}
