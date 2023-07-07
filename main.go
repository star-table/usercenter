package main

import (
	"strconv"

	"github.com/DeanThompson/ginpprof"
	"github.com/dchest/captcha"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/opentracing/opentracing-go"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/star-table/usercenter/core/conf"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/mid"
	"github.com/star-table/usercenter/core/nacos"
	"github.com/star-table/usercenter/core/trace/gin2micro"
	trace "github.com/star-table/usercenter/core/trace/jaeger"
	"github.com/star-table/usercenter/pkg/util/net"
	"github.com/star-table/usercenter/pkg/util/network"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/pkg/util/strs"
	"github.com/star-table/usercenter/service/server"
)

// @title Usercenter API
//// @version v1.0.0
//// @description Movie aggregation search engine.
//// @securityDefinitions.apikey ApiKeyAuth
//// @in header
//// @name Authorization
//// @BasePath /
func main() {
	// 打印程序信息
	logger.Info("redis config: " + strs.ObjectToString(conf.Cfg.Redis))

	logger.Info("mysql config: " + strs.ObjectToString(conf.Cfg.Mysql))

	port := conf.Cfg.Server.Port
	host := conf.Cfg.Server.Host
	if host == "" {
		host = net.GetIP()
	}
	env := conf.GetEnv()

	applicationName := conf.Cfg.Application.Name

	r := gin.New()
	//r.Use(
	//	gin2micro.TracerWrapper,
	//	server.GlsHandler,
	//	server.LoggerHandler,
	//	server.RecoverHandler,
	//)

	// Metrics
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/prometheus")
	m.SetSlowTime(5)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10} used to p95, p99
	m.SetDuration([]float64{0.01, 0.05, 0.1, 0.2, 0.5, 1, 5})
	m.Use(r)

	sentryConfig := conf.Cfg.Sentry
	sentryDsn := ""
	if sentryConfig != nil {
		sentryDsn = sentryConfig.Dsn
	}

	if conf.Cfg.Jaeger != nil {
		jaegerConf := conf.Cfg.Jaeger
		t, io, err := trace.NewTracer(jaegerConf.TraceService, jaegerConf.UdpAddress)
		if err != nil {
			logger.InfoF("err %v", err)
		}
		defer func() {
			if err := io.Close(); err != nil {
				logger.ErrorF("err %v", err)
			}
		}()
		opentracing.SetGlobalTracer(t)

		r.Use(gin2micro.TracerWrapper)
	}

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(mid.SentryMiddleware(applicationName, env, sentryDsn))
	r.Use(mid.StartTrace())
	r.Use(mid.GinContextToContextMiddleware())
	//r.Use(mid.CorsMiddleware())
	r.Use(mid.AuthMiddleware())
	r.Use(mid.TracingHandler())
	captcha.SetCustomStore(&server.RedisCache{})

	//健康检查
	r.GET("/usercenter/health", server.HeartbeatHandler())

	//1.获取验证码
	r.POST("/usercenter/api/task/captcha", server.CaptchaGetHandler())

	//2.获取验证码图片
	r.GET("/usercenter/api/task/captcha/:captchaId", server.CaptchaShowHandler())

	s := server.New(r)
	s.Init()

	if env != consts.RunEnvNull {
		logger.Info("开启pprof监控")
		ginpprof.Wrap(r)
	}

	if ok, _ := slice.Contain([]string{consts.RunEnvBjxLocal, consts.RunEnvBjxTest, consts.RunEnvBjxProd, consts.RunEnvBjxFuseK8s}, env); !ok {
		suc, err := nacos.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          host,
			Port:        uint64(port),
			ServiceName: applicationName,
			GroupName:   conf.Cfg.Discovery.GroupName,
			ClusterName: conf.Cfg.Discovery.ClusterName,
			Weight:      conf.Cfg.Discovery.Weight,
			Enable:      conf.Cfg.Discovery.Enable,
			Healthy:     conf.Cfg.Discovery.Healthy,
			Ephemeral:   conf.Cfg.Discovery.Ephemeral,
			Metadata: map[string]string{
				"kind":    "http",
				"version": "",
			},
		})
		if err != nil {
			logger.Fatal(err)
			return
		}
		if !suc {
			logger.Fatal("服务注册失败")
			return
		}
	}

	logger.InfoF("RUN_ENV:%s, connect to http://%s:%d/ for %s service", env, network.GetIntranetIp(), port, applicationName)
	s.Run(":" + strconv.Itoa(port))
}
