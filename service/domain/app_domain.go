package domain

import (
	"strconv"

	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/nacos"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/service/model/resp"
)

// 调用远程服务获取应用包列表
func GetAppPackageList(orgId int64) ([]resp.AppPackageData, errs.SystemErrorInfo) {
	msg, code, err := nacos.DoGet("app", "app/inner/api/v1/app/packages/get-apppkg-list", map[string]interface{}{"orgId": strconv.FormatInt(orgId, 10)})
	if err != nil {
		logger.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.SystemError, err)
	}
	if code != 200 {
		logger.Error("应用包列表查询接口请求错误")
		return nil, errs.BuildSystemErrorInfo(errs.SystemError)
	}
	var appReq resp.AppPackageListReq
	if msg != "" {
		jsonErr := json.FromJson(msg, &appReq)
		if jsonErr != nil {
			logger.Error(jsonErr)
			return nil, errs.JSONConvertError
		}
	}
	if appReq.Code != 0 {
		logger.Error("应用包列表查询接口请求错误")
		return nil, errs.BuildSystemErrorInfo(errs.SystemError)
	}
	return appReq.Data, nil
}

// 调用远程服务获取app列表
func GetAppList(orgId int64) ([]resp.AppData, errs.SystemErrorInfo) {
	msg, code, err := nacos.DoGet("app", "app/inner/api/v1/apps/get-app-list", map[string]interface{}{"orgId": strconv.FormatInt(orgId, 10)})
	if err != nil {
		logger.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.SystemError, err)
	}
	if code != 200 {
		logger.Error("应用列表查询接口请求错误")
		return nil, errs.BuildSystemErrorInfo(errs.SystemError)
	}
	var appReq resp.AppListReq
	if msg != "" {
		_ = json.FromJson(msg, &appReq)
	}
	if !appReq.Success {
		logger.Error("应用列表查询接口请求错误")
		return nil, errs.BuildSystemErrorInfo(errs.SystemError)
	}
	return appReq.Data, nil
}
