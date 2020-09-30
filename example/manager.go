package main

import (
	"fmt"
	"github.com/gohouse/crontab"
	"github.com/gohouse/golib/date"
	"log"
	"time"
)

func main() {
	var job = crontab.NewTaskManager()

	cron := crontab.NewCronTab(crontab.CT_Second).SetSecond(3)
	cron2 := crontab.NewCronTab(crontab.CT_Second).SetSecond(3)
	job.Add("xxx", cron, teststr)
	job.Add("xxx222", cron2, teststrs)

	log.Println("start...")
	job.Start()
	//go func() {
	//	time.Sleep(10*time.Second)
	//	job.Stop()
	//}()
	job.Wait()
}

func teststr(args ...interface{}) {
	fmt.Println("xxx: ", time.Now().Format(date.DateTimeFormat))
}
func teststrs(args ...interface{}) {
	fmt.Println("xxx222: ", time.Now().Format(date.DateTimeFormat))
}
