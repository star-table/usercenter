module github.com/star-table/usercenter

go 1.13

require (
	github.com/DeanThompson/ginpprof v0.0.0-20190408063150-3be636683586
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869
	github.com/codahale/hdrhistogram v0.9.0 // indirect
	github.com/dchest/captcha v0.0.0-20200903113550-03f5f0333e1f
	github.com/disintegration/imaging v1.6.2
	github.com/emitter-io/go/v2 v2.0.9
	github.com/getsentry/sentry-go v0.11.0
	github.com/gin-contrib/gzip v0.0.1
	github.com/gin-gonic/gin v1.7.4
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.3 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/gopherjs/gopherjs v0.0.0-20190430165422-3e4dfb77656c // indirect
	github.com/json-iterator/go v1.1.12
	github.com/jtolds/gls v4.20.0+incompatible
	github.com/magiconair/properties v1.8.5
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/martinlindhe/base36 v1.1.0
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/logger/zap/v2 v2.9.1
	github.com/modern-go/reflect2 v1.0.2
	github.com/mozillazg/go-pinyin v0.15.0
	github.com/nacos-group/nacos-sdk-go v1.0.8
	github.com/opentracing/opentracing-go v1.1.0
	github.com/penglongli/gin-metrics v0.1.10
	github.com/pkg/errors v0.9.1
	github.com/polaris-team/dingtalk-sdk-golang v0.0.9
	github.com/qustavo/sqlhooks/v2 v2.1.0
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/cast v1.3.1
	github.com/spf13/viper v1.8.1
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.7.0
	github.com/tealeg/xlsx v1.0.5
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	github.com/yoyofxteam/nacos-viper-remote v0.4.0
	github.com/zeta-io/ginx v1.0.1
	github.com/zeta-io/zeta v1.0.3
	go.uber.org/zap v1.17.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/net v0.0.0-20210726213435-c6fcb2dbf985 // indirect
	golang.org/x/tools v0.1.5 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/fatih/set.v0 v0.2.1
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	gotest.tools v2.2.0+incompatible
	upper.io/db.v3 v3.7.1+incompatible
)

replace upper.io/db.v3 v3.7.1+incompatible => github.com/star-table/db v0.3.75-0.20230707012646-28b2e2303a74
