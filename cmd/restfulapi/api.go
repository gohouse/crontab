package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gohouse/crontab"
	"github.com/gohouse/date"
	"github.com/gohouse/t"
	"log"
	"net/http"
	"time"
)

var tm *crontab.TaskManager
//var logfile = "cron.log"

func init() {
	tm = crontab.NewTaskManager()
	//cron1 := crontab.NewCronTab(crontab.CT_Second).SetSecond(5)
	//cron2 := crontab.NewCronTab(crontab.CT_Second).SetSecond(10)
	//tm.Add(uuid.NewV4().String(), cron1, func(args ...interface{}) {
	//	file.FilePutContents(logfile, []byte(fmt.Sprintf("%s cron1的测试\n", time.Now().Format(date.DateTimeFormat))))
	//}).Add(uuid.NewV4().String(), cron2, func(args ...interface{}) {
	//	file.FilePutContents(logfile, []byte(fmt.Sprintf("%s cron2的测试\n", time.Now().Format(date.DateTimeFormat))))
	//})
}
func main() {
	route := gin.Default()
	route.GET("/", func(c *gin.Context) {
		tmp := fmt.Sprintf(`<h1>计划任务说明</h1>可以通过api添加指定秒数的周期任务,如:<br> 
启动给你一个周期为3s的测试计划任务: <a target='_blank' href='/new/3'>开始任务:/new/3</a><br>
启动给你一个周期为5s的测试计划任务: <a target='_blank' href='/new/5'>开始任务:/new/5</a><br>
停止任务 api: /stop/{id}<br>
启动任务 api: /start/{id}<br>
移出任务 api: /remove/{id}<br>
生成一个测试任务 api: /new/{seconds}<br>
任务列表 api: <a target='_blank' href='/taskList'>/taskList</a><br>
`)
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, tmp)
	})
	route.GET("/taskList", taskList)
	route.GET("/new/:step", refresh)
	route.GET("/stop/:name", stop)
	route.GET("/start/:name", start)
	route.GET("/remove/:name", remove)
	//route.GET("/cronlog", cronlog)

	log.Println("view: http://localhost:8988")
	log.Fatal(route.Run(":8988"))
}

func start(c *gin.Context) {
	var name = c.Param("name")
	tm.Start(name)
	c.String(http.StatusOK, "启动:"+name)
}
func stop(c *gin.Context) {
	var name = c.Param("name")
	tm.Stop(name)
	c.String(http.StatusOK, "停止:"+name)
}
func remove(c *gin.Context) {
	var name = c.Param("name")
	tm.Remove(name)
	c.String(http.StatusOK, "移出:"+name)
}
func taskList(c *gin.Context) {
	var tmp = "<h1>任务列表</h1>"
	tm.Range(func(key, value interface{}) bool {
		val := value.(*crontab.TaskObject)
		var taskStatus = "已停止"
		if val.IsRunning() {
			taskStatus = "运行中"
			tmp += fmt.Sprintf("任务id: %s; 状态: <b style='color:green'>%s</b>; 操作: <a target='_blank' href='/stop/%v'>停止:/stop/%v</a><br>",
				key, taskStatus, key, key)
		} else {
			tmp += fmt.Sprintf("任务id: %s; 状态: <b style='color:red'>%s</b>; 操作: <a target='_blank' href='/start/%v'>启动:/stop/%v</a> <a target='_blank' href='/remove/%v'>移出:/remove/%v</a><br>",
				key, taskStatus, key, key, key, key)
		}
		return true
	})
	if tmp=="<h1>任务列表</h1>" {
		tmp += fmt.Sprintf(`暂无任务,可以通过api添加指定秒数的周期任务,如:<br> 
启动给你一个周期为3s的测试计划任务: <a target='_blank' href='/new/3'>开始任务:/new/3</a><br>
启动给你一个周期为5s的测试计划任务: <a target='_blank' href='/new/5'>开始任务:/new/5</a><br>
停止任务 api: /stop/{id}<br>
启动任务 api: /start/{id}<br>
移出任务 api: /remove/{id}<br>
生成一个测试任务 api: /new/{seconds}<br>
任务列表 api: <a target='_blank' href='/taskList'>/taskList</a><br>
`)
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, tmp)
}

func refresh(c *gin.Context) {
	var step = c.Param("step")
	var taskName = fmt.Sprintf("test_%vs", step)
	var testcron = crontab.NewCronTab(crontab.CT_Second).SetSecond(t.New(step).Int()).RunOnceFirst()
	var callback = func(args ...interface{}) {
		log.Printf("%v: 每 %vs 刷新一次\n", time.Now().Format(date.DateTimeFormat), args[0])
	}
	tm.Add(taskName, testcron, callback, step)
	tm.Start(taskName)

	w := c.Writer
	f := w.(http.Flusher)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	fmt.Fprintf(w, "新增加任务: %s, 任务周期为: %ss, 停止本任务,请调用stop api: <a href='/stop/%v'>/stop/%v</a> <br>", taskName, step, taskName, taskName)
	f.Flush()
}
