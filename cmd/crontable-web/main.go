package main

import (
	"github.com/gohouse/crontab"
	"github.com/gohouse/crontab/adapter/fileLog"
	"github.com/gohouse/crontab/client"
	"log"
)

func main() {
	var port = ":8200"
	// 日志文件
	var logfile = "crontab.log"

	// 初始化计划任务管理器
	// 这里使用了文件记录日志, 提供了logger接口, 可以自由扩充记录日志到数据库等其他地方
	tm := crontab.NewTaskManager(crontab.Logger(fileLog.NewFileLog(logfile)))

	// 加入任务列表
	//TaskInit(tm)
	tm.AddGroup(TaskInit)

	// 开启 restful api 服务
	log.Fatal(client.Run(tm, port))
}

func TaskInit(tm *crontab.TaskManager) {
	// test 每10s执行一次
	tm.Add("每10s执行一次", crontab.NewCronTab(crontab.CT_Second).SetSecond(10), Test)

	// 每天执行一次  statistic_of_per_day
	tm.Add("每天0时0分0秒执行的任务", crontab.NewCronTab(crontab.CT_Day).SetDay(1), Test)

	// 30分钟执行一次
	tm.Add("30分钟执行一次",
		crontab.NewCronTab(crontab.CT_Minute).SetMinute(30).
		RunOnceFirst(false),	// 这一步操作是移除默认先执行一次,而是从30分钟后的 0s 开始周期执行第一次
		Test)
}

func Test(args ...interface{}) {
	// todo 这里就是你想干的事
}
