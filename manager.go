package crontab

import (
	"context"
	"sync"
)

type TaskObject struct {
	cron     *CronTab
	callback HandleFunc
	args     []interface{}
	cancel   context.CancelFunc
}

func (so *TaskObject) start() {
	if so.cron.running == true {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	so.cron.Context = ctx
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

type TaskManager struct {
	*sync.Map
	wg  *sync.WaitGroup
	ctx context.Context
}

func NewTaskManager() *TaskManager {
	return &TaskManager{&sync.Map{}, &sync.WaitGroup{}, context.Background()}
}

func (job *TaskManager) Add(key string, cron *CronTab, callback HandleFunc, args ...interface{}) *TaskManager {
	var so = TaskObject{
		cron:     cron,
		callback: callback,
		args:     args,
	}
	job.Store(key, &so)
	return job
}

func (job *TaskManager) Start(keys ...string) {
	if len(keys) > 0 {
		if r, ok := job.Load(keys[0]); ok {
			var so = r.(*TaskObject)
			job.wg.Add(1)
			go so.start()
			job.wg.Done()
		}
	} else {
		job.Range(func(_, value interface{}) bool {
			job.wg.Add(1)
			var so = value.(*TaskObject)
			go so.start()
			job.wg.Done()
			return true
		})
	}
}

func (job *TaskManager) Wait() {
	job.wg.Wait()
}

func (job *TaskManager) Stop(keys ...string) {
	if len(keys) > 0 {
		if r, ok := job.Load(keys[0]); ok {
			var so = r.(*TaskObject)
			so.stop()
		}
	} else {
		job.Range(func(_, value interface{}) bool {
			var so = value.(*TaskObject)
			so.stop()
			return true
		})
	}
}

func (job *TaskManager) Remove(keys ...string) {
	if len(keys) > 0 {
		job.Stop(keys[0])
		job.Delete(keys[0])
	} else {
		job.Stop()
		job = NewTaskManager()
	}
}
