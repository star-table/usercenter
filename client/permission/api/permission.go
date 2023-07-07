package api

import (
	"errors"

	"github.com/star-table/usercenter/client/permission/model/req"
	"github.com/star-table/usercenter/client/permission/model/resp"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/nacos"
	"github.com/star-table/usercenter/pkg/util/json"
)

type permissionClient byte

var PermissionClient permissionClient

func (*permissionClient) InitDefaultManageGroup(req req.InitDefaultManageGroupReq) (bool, error) {
	respBody, _, err := nacos.DoPost(consts.ServiceNamePermission, "permission/inner/api/v1/manage-group/init", nil, json.ToJsonIgnoreError(req))
	if err != nil {
		logger.Error(err)
		return false, err
	}
	result := resp.Result{}
	json.FromJsonIgnoreError(respBody, &result)

	if result.Code != 0 {
		logger.ErrorF("InitDefaultManageGroup resp err: %v， res data: %s", err, respBody)
		return false, errors.New(result.Message)
	}

	return result.Data == true, nil
}

func (*permissionClient) InitAppPermission(req req.InitAppPermissionReq) (bool, error) {
	respBody, _, err := nacos.DoPost(consts.ServiceNamePermission, "permission/inner/api/v1/app-permission/init", nil, json.ToJsonIgnoreError(req))
	if err != nil {
		logger.Error(err)
		return false, err
	}
	result := resp.Result{}
	json.FromJsonIgnoreError(respBody, &result)

	if result.Code != 0 {
		logger.ErrorF("InitDefaultManageGroup resp err: %v， res data: %s", err, respBody)
		return false, errors.New(result.Message)
	}

	return result.Data == true, nil
}
