# crontab
a simple and powerful crontab written in golang with web page management  
golang实现的简单便捷的计划任务管理系统, 自带 web 界面,方便的管理多个任务  
支持 秒,分,时,日,月,周  

## 管理界面
![](example/demo.jpg)

## install
- go.mod
```shell script
require github.com/gohouse/crontab master
````

## web管理简单用例
```go
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
	TaskInit(tm)

	// 开启 restful api 服务
	log.Fatal(client.Run(tm, port))
}

func TaskInit(tm *crontab.TaskManager) {
	// test 每3s执行一次
	tm.Add("每3s执行一次", crontab.NewCronTab(crontab.CT_Second).SetSecond(3), Test)

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
```
访问 http://localhost:8200 即可方便的查看管理计划任务了

## 非web管理用例
### 1. 执行简单的任务
```go
package main

import (
	"github.com/gohouse/crontab"
	"log"
)

func main() {
	crontab.NewCronTab(crontab.CT_Second).
		SetSecond(3).
		Run(func(args ...interface{}) {
			log.Println("每 3s 会执行一次本操作")
		})
}
```
输出
```go
2020/01/23 20:26:29 每 3s 会执行一次本操作
2020/01/23 20:26:32 每 3s 会执行一次本操作
2020/01/23 20:26:35 每 3s 会执行一次本操作
```

### 2. 执行多个任务
```go
package main

import (
	"fmt"
	"github.com/gohouse/crontab"
	"github.com/gohouse/date"
	"log"
	"time"
)

func main() {
	var job = crontab.NewTaskManager()
	cron := crontab.NewCronTab(crontab.CT_Second).SetSecond(1)
	cron2 := crontab.NewCronTab(crontab.CT_Second).SetSecond(1)
	job.Add("xxx", cron, teststr).Add("xxx222", cron2, teststrs)

	log.Println("start...")
	job.Start()
	job.Wait()
}

func teststr(args ...interface{}) {
	fmt.Println("xxx: ", time.Now().Format(date.DateTimeFormat))
}
func teststrs(args ...interface{}) {
	fmt.Println("xxx222: ", time.Now().Format(date.DateTimeFormat))
}
```

### restful api 
参考 `cmd/restfulapi/api.go`  
`GET /start/{id}`  
`GET /stop/{id}`  
`GET /remove/{id}`  
`GET /tasklist`  

### api参考
```shell script
# 每3s执行一次
crontab.NewCronTab(crontab.CT_Second).SetSecond(3)
# 每分钟的第5s执行一次
crontab.NewCronTab(crontab.CT_Minute).SetMinute(1).SetSecond(5)
# 每2小时的第0分5s执行一次,缺省分钟则默认为0,下同
crontab.NewCronTab(crontab.CT_Hour).SetHour(2).SetSecond(5)
# 每3天的0点0分5s执行一次
crontab.NewCronTab(crontab.CT_Day).SetDay(3).SetSecond(5)
# 每月1号的03点05分0s执行一次,缺省日期为1号
crontab.NewCronTab(crontab.CT_Month).SetMonth(1).SetHour(3).SetMinute(5)
# 每周周日的0点5分0s执行一次
crontab.NewCronTab(crontab.CT_Week).SetWeek(time.Sunday).SetMinute(5)
```
所有计划任务再运行时,都会优先执行一次,如果不想先执行一次,则可以调用`RunOnceFirst(false)`即可,如
```shell script
crontab.NewCronTab(crontab.CT_Second).SetSecond(3).RunOnceFirst(false)
```