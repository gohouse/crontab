package client

import (
	"github.com/gin-gonic/gin"
	"github.com/gohouse/crontab"
	"github.com/gohouse/crontab/client/client_web"
	"github.com/gohouse/t"
	"net/http"
	"strings"
)

var tm *crontab.TaskManager
var htmlRaw = client_web.LoadTemplate()

func Run(ctm *crontab.TaskManager, port string) error {
	tm = ctm
	// 启动web服务
	route := gin.Default()
	routeInit(route)
	return route.Run(port)
}

func routeInit(route *gin.Engine) {
	route.GET("/", index)
	route.GET("/tasklist", taskList)
	route.GET("/new/:step", refresh)
	route.GET("/stop/:name", stop)
	route.GET("/start", startAll)
	route.GET("/start/:name", start)
	route.GET("/remove/:name", remove)
	route.GET("/log", logInfo)
}

func index(c *gin.Context) {
	//c.Header("Content-Type", "text/html; charset=utf-8")
	//c.String(http.StatusOK, notice)

	//c.HTML(http.StatusOK, "index.html", nil)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, htmlRaw)
}

func start(c *gin.Context) {
	var name = c.Param("name")
	tm.Start(name)
	c.String(http.StatusOK, "启动:"+name)
}

func startAll(c *gin.Context) {
	tm.Start()
	c.String(http.StatusOK, "启动所有任务")
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
	var result = []map[string]interface{}{}
	tm.Range(func(key, value interface{}) bool {
		val := value.(*crontab.TaskObject)
		var taskStatus = "已停止"
		if val.IsRunning() {
			taskStatus = "运行中"
		}
		result = append(result, map[string]interface{}{
			"status":taskStatus,
			"id":key,
			"title":val.Title(),
		})
		return true
	})
	//c.Header("Content-Type", "text/html; charset=utf-8")
	c.JSON(http.StatusOK, result)
}

func refresh(c *gin.Context) {
	var step = c.Param("step")

	tm.Add("周期"+step+"s测试任务",
		crontab.NewCronTab(crontab.CT_Second).SetSecond(t.New(step).Int()).RunOnceFirst(),
		Test)
}

func logInfo(c *gin.Context) {
	var limit int64 = 20
	if r,ok:=c.GetQuery("limit");ok{
		limit = t.New(r).Int64()
	}
	c.JSON(http.StatusOK, strings.Split(tm.LogInfo(limit),"\n"))
}

func Test(args ...interface{})  {
	// todo
}
