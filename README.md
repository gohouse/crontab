# crontab
a simple and powerful crontab write in golang

## install
- go.mod
```shell script
require github.com/gohouse/crontab master
```

## 简单用例
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