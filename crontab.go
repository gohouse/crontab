package crontab

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync/atomic"
	"time"
)

type CronType int

const (
	CT_Month CronType = iota + 1
	CT_Week
	CT_Day
	CT_Hour
	CT_Minute
	CT_Second
)

type CronValue struct {
	Month  int
	Week   time.Weekday
	Day    int
	Hour   int
	Minute int
	Second int
}

type HandleFunc func(args ...interface{})

type CronTab struct {
	ctx context.Context
	CronType
	cronv        CronValue
	runOnceFirst bool
	running      bool
	opt          *Options
	runTimes     int64
}

func NewCronTab(cron CronType, opts ...OptionHandleFunc) *CronTab {
	var opt = &Options{}
	for _, item := range opts {
		item(opt)
	}
	return &CronTab{
		ctx:      context.TODO(),
		CronType: cron,
		cronv: CronValue{
			Month: 1,
			Day:   1,
		},
		runOnceFirst: true,
		opt:          opt,
	}
}

// RunOnceFirst 先运行一次
func (ct *CronTab) RunOnceFirst(b ...bool) *CronTab {
	if len(b) > 0 {
		ct.runOnceFirst = b[0]
	} else {
		ct.runOnceFirst = true
	}

	return ct
}

func (ct *CronTab) SetMonth(arg int) *CronTab {
	if arg < 1 {
		panic("arg must > 0")
	}
	ct.cronv.Month = arg
	return ct
}

func (ct *CronTab) SetWeek(arg time.Weekday) *CronTab {
	ct.cronv.Week = arg
	return ct
}

func (ct *CronTab) SetDay(arg int) *CronTab {
	if arg < 0 {
		panic("arg must >= 0")
	}
	ct.cronv.Day = arg
	return ct
}

func (ct *CronTab) SetHour(arg int) *CronTab {
	if arg < 0 {
		panic("arg must >= 0")
	}
	ct.cronv.Hour = arg
	return ct
}

func (ct *CronTab) SetMinute(arg int) *CronTab {
	if arg < 0 {
		panic("arg must >= 0")
	}
	ct.cronv.Minute = arg
	return ct
}

func (ct *CronTab) SetSecond(arg int) *CronTab {
	if arg < 0 {
		panic("arg must >= 0")
	}
	ct.cronv.Second = arg
	return ct
}

func (ct *CronTab) IsRunning() bool {
	return ct.running
}

func (ct *CronTab) RunTimes() int64 {
	return ct.runTimes
}

func (ct *CronTab) Run(h HandleFunc, args ...interface{}) {
	if ct.opt.logger == nil {
		ct.opt.logger = logrus.New()
	}
	ct.running = true
	if ct.runOnceFirst {
		go h(args...)
		atomic.AddInt64(&ct.runTimes, 1)
		ct.opt.logger.Infof("第%v次执行任务:%v", ct.runTimes, args)
	}
	for {
		now := time.Now()
		next := ct._run(now)
		t := time.NewTimer(next.Sub(now))
		defer t.Stop()
		select {
		case <-ct.ctx.Done():
			ct.running = false
			//log.Println("done ...")
			return
		case <-t.C:
			//以下为定时执行的操作
			go h(args...)
			atomic.AddInt64(&ct.runTimes, 1)
			ct.opt.logger.Infof("第%v次执行任务:%v", ct.runTimes, args)
		}
	}
}

func (ct *CronTab) _run(now time.Time) time.Time {
	switch ct.CronType {
	case CT_Month:
		if ct.cronv.Month == 0 {
			panic("Month must > 0")
		}
		next := now.AddDate(0, ct.cronv.Month, 0)
		return time.Date(next.Year(), next.Month(), ct.cronv.Day, ct.cronv.Hour, ct.cronv.Minute, ct.cronv.Second, 0, next.Location())
	case CT_Week:
		var days = time.Saturday - now.Weekday() + ct.cronv.Week + 1
		next := now.AddDate(0, int(days), 0)
		return time.Date(next.Year(), next.Month(), next.Day(), ct.cronv.Hour, ct.cronv.Minute, ct.cronv.Second, 0, next.Location())
	case CT_Day:
		if ct.cronv.Day == 0 {
			panic("Day must > 0")
		}
		next := now.AddDate(0, 0, ct.cronv.Day)
		return time.Date(next.Year(), next.Month(), next.Day(), ct.cronv.Hour, ct.cronv.Minute, ct.cronv.Second, 0, next.Location())
	case CT_Hour:
		if ct.cronv.Hour == 0 {
			panic("Hour must > 0")
		}
		next := now.Add(time.Hour * time.Duration(ct.cronv.Hour))
		return time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), ct.cronv.Minute, ct.cronv.Second, 0, next.Location())
	case CT_Minute:
		if ct.cronv.Minute == 0 {
			panic("Minute must > 0")
		}
		next := now.Add(time.Minute * time.Duration(ct.cronv.Minute))
		return time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), ct.cronv.Second, 0, next.Location())
	case CT_Second:
		if ct.cronv.Second == 0 {
			panic("Second must > 0")
		}
		next := now.Add(time.Second * time.Duration(ct.cronv.Second))
		return time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), next.Second(), 0, next.Location())
	}
	return time.Now()
}
