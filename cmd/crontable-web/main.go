package main

import (
	"encoding/json"
	"fmt"
	"github.com/gohouse/crontab"
	"github.com/gohouse/crontab/client"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

type logFormater struct {}
func (logFormater) Format(entry *logrus.Entry) ([]byte, error) {
	var marshal []byte
	if len(entry.Data) > 0 {
		marshal, _ = json.Marshal(entry.Data)
	}
	res := fmt.Sprintf("[%s] %s %s %s\n", entry.Level.String(), entry.Time.Format(time.RFC3339), entry.Message, marshal)
	return []byte(res),nil
}
func main() {
	var port = ":8200"
	// 日志
	logger := logrus.New()

	// 如果使用日志文件
	var logfile = "crontab.log"
	f, _ := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0766)
	logger.SetOutput(f)
	logger.SetFormatter(logFormater{})

	// 如果输出到控制台,则使用 os.Stdout 即可
	//logger.SetOutput(os.Stdout)

	// 初始化计划任务管理器
	// 这里使用了文件记录日志, 提供了logger接口, 可以自由扩充记录日志到数据库等其他地方
	tm := crontab.NewTaskManager(crontab.Logger(logger))

	// 加入任务列表
	//TaskInit(tm)
	tm.AddGroup(TaskInit)

	// 开启 restful api 服务
	log.Fatal(client.Run(tm, port, logfile))
}

func TaskInit(tm *crontab.TaskManager) {
	// test 每10s执行一次
	tm.Add("每10s执行一次", crontab.NewCronTab(crontab.CT_Second).SetSecond(10), Test)

	// 每天执行一次  statistic_of_per_day
	tm.Add("每天0时0分0秒执行的任务", crontab.NewCronTab(crontab.CT_Day).SetDay(1), Test)

	// 30分钟执行一次
	tm.Add("30分钟执行一次",
		crontab.NewCronTab(crontab.CT_Minute).SetMinute(30).
			RunOnceFirst(false), // 这一步操作是移除默认先执行一次,而是从30分钟后的 0s 开始周期执行第一次
		Test)
}

func Test(args ...interface{}) {
	// todo 这里就是你想干的事
}
