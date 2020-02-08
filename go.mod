module github.com/gohouse/crontab

go 1.13

require (
	github.com/gin-gonic/gin v1.5.0
	github.com/gohouse/date v0.0.0-20191203070644-0e6619197ef8
	github.com/gohouse/file v0.0.0-20200205051838-0d350a0b6f2b
	github.com/gohouse/random v0.0.0-20200102081411-fd6c47e80d2c
	github.com/gohouse/t v0.0.5
)

replace github.com/gohouse/file => ../file
