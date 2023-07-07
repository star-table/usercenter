package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/logger"
)

func LoggerHandler(c *gin.Context) {
	// Start timer
	start := time.Now()

	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	method := c.Request.Method
	// Process request
	c.Next()

	// Stop timer
	end := time.Now()
	latency := end.Sub(start)
	statusCode := c.Writer.Status()
	ecode := c.GetInt("500")
	clientIP := c.ClientIP()
	if raw != "" {
		path = path + "?" + raw
	}
	logger.DebugF("METHOD:%s | PATH:%s | CODE:%d | IP:%s | TIME:%d | ECODE:%d \n", method, path, statusCode, clientIP, latency/time.Millisecond, ecode)
}

func RecoverHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			httprequest, _ := httputil.DumpRequest(c.Request, false)
			pnc := fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s\n%s", time.Now().Format("2006-01-02 15:04:05"), string(httprequest), err, buf)
			fmt.Print(pnc)
			c.AbortWithStatus(500)
		}
	}()
	c.Next()
}

func SwaggerHandler() gin.HandlerFunc {
	fs := http.FileServer(http.Dir("./../swagger"))
	h := http.StripPrefix("/swagger/", fs)
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func HeartbeatHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, "ok")
	}
}
