package main

import (
	"github.com/gohouse/crontab"
	"log"
)

func main() {
	crontab.NewCronTab(crontab.CT_Second).
		SetSecond(1).
		RunOnceFirst().
		Run(func(args ...interface{}) {
			log.Println("每 3s 会执行一次本操作")
		})
}
