module github.com/gohouse/crontab

go 1.13

require (
	github.com/gin-gonic/gin v1.5.0
	github.com/gohouse/date v0.0.0-20191203070644-0e6619197ef8
	github.com/gohouse/file v0.0.0-20191230075216-864710533ca2
	github.com/gohouse/t v0.0.5
	github.com/satori/go.uuid v1.2.0
)

replace github.com/gohouse/file => ../file
