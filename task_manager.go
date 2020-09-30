package crontab

import (
	"context"
	"github.com/gohouse/golib/t"
	"github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
)

var incId int64

func GetId() int64 { return atomic.AddInt64(&incId, 1) }

type TaskManager struct {
	store *sync.Map
	wg    *sync.WaitGroup
	ctx   context.Context
	opt   *Options
	done  chan struct{}
}

func NewTaskManager(opts ...OptionHandleFunc) *TaskManager {
	var opt = &Options{}
	for _, item := range opts {
		item(opt)
	}
	if opt.logger == nil {
		opt.logger = logrus.New()
	}
	return newTaskManager(opt)
}

func newTaskManager(opt *Options) *TaskManager {
	return &TaskManager{&sync.Map{}, &sync.WaitGroup{}, context.Background(), opt, make(chan struct{})}
}

func (job *TaskManager) Add(title string, cron *CronTab, callback HandleFunc, args ...interface{}) string {
	var taskId = t.New(GetId()).String()
	args = append(args, taskId, "-", title)
	cron.opt = job.opt
	var so = TaskObject{
		cron:     cron,
		callback: callback,
		args:     args,
		title:    title,
		taskId:   taskId,
	}
	job.store.Store(taskId, &so)
	job.opt.logger.Infof("添加任务:%s - %s", taskId, title)
	return taskId
}

func (job *TaskManager) AddGroup(tl func(*TaskManager)) {
	tl(job)
}

func (job *TaskManager) Start(keys ...string) {
	if len(keys) > 0 {
		if r, ok := job.store.Load(keys[0]); ok {
			var so = r.(*TaskObject)
			job.wg.Add(1)
			go so.start()
			job.wg.Done()
			job.opt.logger.Infof("开始任务:%s - %s", so.taskId, so.title)
		}
	} else {
		job.store.Range(func(key, value interface{}) bool {
			job.wg.Add(1)
			var so = value.(*TaskObject)
			go so.start()
			job.wg.Done()
			job.opt.logger.Infof("开始任务:%s - %s", so.taskId, so.title)
			return true
		})
	}
}

func (job *TaskManager) Wait() {
	job.wg.Wait()
	select {}
}

func (job *TaskManager) Stop(keys ...string) {
	if len(keys) > 0 {
		if r, ok := job.store.Load(keys[0]); ok {
			var so = r.(*TaskObject)
			if so.IsRunning() {
				so.stop()
				job.opt.logger.Infof("停止任务:%s - %s", so.taskId, so.title)
			}
		}
	} else {
		job.store.Range(func(key, value interface{}) bool {
			var so = value.(*TaskObject)
			so.stop()
			if so.IsRunning() {
				so.stop()
				job.opt.logger.Infof("停止任务:%s - %s", so.taskId, so.title)
			}
			return true
		})
	}
	//// 判断是否还有任务
	//var jobs int
	//job.store.Range(func(key, value interface{}) bool {
	//	jobs++
	//	return true
	//})
	//if jobs == 0 {
	//	job.done <- struct{}{}
	//}
}

func (job *TaskManager) Remove(keys ...string) {
	if len(keys) > 0 {
		job.Stop(keys[0])
		job.store.Delete(keys[0])
		job.opt.logger.Infof("删除任务:%s", keys[0])
	} else {
		job.opt.logger.Infof("删除所有任务")
		job.Stop()
		*job = *newTaskManager(job.opt)
	}
}

func (job *TaskManager) Range(f func(key, value interface{}) bool) {
	job.store.Range(f)
}
