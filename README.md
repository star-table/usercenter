# 无码平台-用户中心 etc
用户中心提供的是 RESTful 接口，接口路由的声明可以在 service/server/server.go 文件中查看。

## 接口文档
接口文档使用 swagger 生成，生成方式可以[参考这里](https://github.com/swaggo/gin-swagger)

下载安装 `go get -u github.com/swaggo/swag/cmd/swag`。然后就能使用 `swag init`

而在 usercenter 项目中，我们一般在本地启动服务（`go run main.go`）后，通过浏览器访问 http://127.0.0.1:15001/usercenter/swagger/index.html，
可以访问到 swagger 文档页

## duty
* Nico（刘千源）、樊宇、王世昌、苏汉宇

## 开发规范
### 路由声明
restful 接口路由中的路径使用中横线分隔的风格，形如 `needLoginGroup.POST("/org/generate-api-key", api.Org.GenerateApiKey)` 中的 `/org/generate-api-key`

### model 结构体声明
* 接口入参：接口的入参结构体**必须**在 `service/model/req` 目录下声明。
* 接口出参：接口的出参结构体**必须**在 `service/model/resp` 目录下声明。
* 数据库 model：表相关的 model 结构体声明必须在 `service/model/po` 目录下
* service 和 domain：业务内需要定义的结构体可以定义在 `service/model/bo` 目录下

### domain 层
该层尽量不要写太多的业务相关的逻辑。业务相关尽量抽象到 service（service/service） 层。另外，domain 层的方法的异常返回**必须**使用 error，而非 `errs.SystemErrorInfo`

### 注释规范
* service 目录下的业务方法，必须要有注释，并且注释格式形如下方：

```
/**
ModifyPositionInfo 修改职级信息
@author WangShiChang
@version v1.0
@date 2020-10-21
*/
```

其中 ModifyPositionInfo  是方法名称，其后是描述信息。`@author` 表示代码作者；`@date` 是日期

### mutation 接口日志
* 当在编写非查询类接口时，service 层中对接控制器的方法中，必须要进行日志记录，主要记录修改行为涉及到的参数，如：

```
// logger.InfoF("[删除职级] -> 参数 orgPositionId: %d", orgPositionId)
// logger.InfoF("[修改职级信息] -> 参数 orgPositionId: %d, reqParam: %s", orgPositionId, json.ToJsonIgnoreError(reqParam))
```

而如果方法是 query 类型，则可以不记录。
