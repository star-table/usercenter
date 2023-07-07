package mid

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jtolds/gls"
	"github.com/opentracing/opentracing-go"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/threadlocal"
	"github.com/star-table/usercenter/core/trace/gin2micro"
	"github.com/star-table/usercenter/pkg/util/network"
	"github.com/star-table/usercenter/pkg/util/strs"
	"github.com/star-table/usercenter/pkg/util/times"
	"github.com/star-table/usercenter/pkg/util/uuid"
	"github.com/star-table/usercenter/service/model/vo"
)

func newTraceId() string {
	return fmt.Sprintf("%s_%s", network.GetIntranetIp(), uuid.NewUuid())
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func isBinaryContent(contentType string) bool {
	return strings.Contains(contentType, "image") || strings.Contains(contentType, "video") ||
		strings.Contains(contentType, "audio")
}

func isMultipart(contentType string) bool {
	return strings.Contains(contentType, "multipart/form-data")
}

func TracingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		openTraceId := ""
		var span opentracing.Span
		traceContext, ok := gin2micro.ContextWithSpan(c)
		if ok && traceContext != nil {
			if v := traceContext.Value(gin2micro.OpenTraceKey); v != nil {
				openTraceId = v.(string)
			}
			span = opentracing.SpanFromContext(traceContext)
		}
		threadlocal.Mgr.SetValues(gls.Values{consts.TraceIdKey: openTraceId, consts.JaegerContextTraceKey: openTraceId, consts.JaegerContextSpanKey: span}, func() {
			c.Next()
		})
	}
}

func StartTrace() gin.HandlerFunc {
	return func(c *gin.Context) {

		traceId := c.Request.Header.Get(consts.TraceIdKey)
		if "" == traceId {
			traceId = newTraceId()
			traceContext, ok := gin2micro.ContextWithSpan(c)
			if ok && traceContext != nil {
				if v := traceContext.Value(gin2micro.OpenTraceKey); v != nil {
					traceId = v.(string)
				}
			}
		}

		contentType := c.ContentType()
		//postForm := ""
		body := ""
		httpContext := vo.HttpContext{
			Request:   c.Request,
			Url:       c.Request.RequestURI,
			Method:    c.Request.Method,
			StartTime: times.GetNowMillisecond(),
			EndTime:   0,
			Ip:        c.ClientIP(),
			TraceId:   traceId,
		}

		userToken := c.GetHeader(consts.IdentityUserHeader)
		orgId := c.GetHeader(consts.IdentityOrgHeader)
		env := c.GetHeader(consts.AppHeaderEnvName)
		plat := c.GetHeader(consts.AppHeaderPlatName)
		proId := c.GetHeader(consts.AppHeaderProName)
		ver := c.GetHeader(consts.AppHeaderVerName)

		if !isBinaryContent(contentType) && !isMultipart(contentType) {
			// 判断不是上传文件等大消息体，记录消息体日志
			//c.Request.ParseForm()
			//postForm = c.Request.PostForm.Encode()
			data, err := c.GetRawData()
			if err != nil {
				logger.Info(err.Error())
			}
			body = string(data)
			// 重新写入body
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		}

		ctx := context.WithValue(c.Request.Context(), consts.TraceIdKey, traceId)
		ctx = context.WithValue(ctx, consts.HttpContextKey, httpContext)
		c.Request = c.Request.WithContext(ctx)

		c.Header(consts.TraceIdKey, traceId)

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		urlCount := strs.Len(httpContext.Url)
		httpContext.EndTime = times.GetNowMillisecond()
		runtime := httpContext.EndTime - httpContext.StartTime
		if urlCount <= 6 || (httpContext.Url[urlCount-6:] != "health" || runtime >= 100) {
			// 仅记录非心跳日志

			httpContext = c.Request.Context().Value(consts.HttpContextKey).(vo.HttpContext)
			httpContext.Status = c.Writer.Status()
			runtimes := strconv.FormatInt(runtime, 10)
			httpStatus := strconv.Itoa(httpContext.Status)
			msg := fmt.Sprintf("request|traceId=%s|start=%s|ip=%s|contentType=%s|method=%s|url=%s|header={\"ver\":\"%s\","+
				"\"user\":\"%s\",\"plat\":\"%s\",\"env\":\"%s\",\"org\":\"%s\",\"pro\":\"%s\"}|body=%s|------response|end=%s|time=%s|status=%s|body=%s|",
				httpContext.TraceId, times.GetDateTimeStrByMillisecond(httpContext.StartTime), httpContext.Ip, contentType,
				httpContext.Method, httpContext.Url, ver, userToken, plat, env, orgId, proId, strings.ReplaceAll(body, "\n", "\\n"), times.GetDateTimeStrByMillisecond(httpContext.EndTime),
				runtimes, httpStatus, blw.body.String())
			//logger.InfoF("traceId=%s,%v,duration:%v,responseCode:%v,path:%v", httpContext.TraceId, msg, runtimes, httpStatus, httpContext.Url)
			//打印要优化，数据太多会刷屏
			logger.InfoF("traceId=%s,%v,duration:%v,responseCode:%v,path:%v", httpContext.TraceId, len(msg), runtimes, httpStatus, httpContext.Url)
		}
	}
}
