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
	route.LoadHTMLGlob("client/client_web/*")
	route.GET("/", index)
	route.GET("/tasklist", taskList)
	route.GET("/new/:step", refresh)
	route.GET("/stop", stop)
	route.GET("/stop/:pkid", stop)
	route.GET("/start", start)
	route.GET("/start/:pkid", start)
	route.GET("/remove", remove)
	route.GET("/remove/:pkid", remove)
	route.GET("/log", logInfo)
}

func index(c *gin.Context) {
	//c.HTML(http.StatusOK, "index.html", nil)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, htmlRaw)
}

func start(c *gin.Context) {
	var pkid = c.Param("pkid")
	if pkid=="" {
		tm.Start()
	} else {
		tm.Start(pkid)
	}
	c.String(http.StatusOK, "启动:"+pkid)
}
func stop(c *gin.Context) {
	var pkid = c.Param("pkid")
	if pkid=="" {
		tm.Stop()
	} else {
		tm.Stop(pkid)
	}
	c.String(http.StatusOK, "停止:"+pkid)
}
func remove(c *gin.Context) {
	var pkid = c.Param("pkid")
	if pkid=="" {
		tm.Remove()
	} else {
		tm.Remove(pkid)
	}
	c.String(http.StatusOK, "删除:"+pkid)
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
