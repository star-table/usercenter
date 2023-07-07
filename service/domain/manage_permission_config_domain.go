package domain

import (
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/service/model/po"
	"upper.io/db.v3"
)

// GetManageGroupConfig 获取管理组
func GetManageGroupConfig() (*po.LcPerManagePermissionConfig, error) {
	var mpc po.LcPerManagePermissionConfig
	dbErr := store.Mysql.SelectOneByCond(consts.TableManagePermissionConfig, db.Cond{
		consts.TcOrgId:   0,
		consts.TcType:    2,
		consts.TcDelFlag: consts.AppIsNoDelete,
	}, &mpc)

	if dbErr != nil {
		return nil, dbErr
	}

	return &mpc, nil
}

// GetManageGroupConfigForPolaris 获取管理组，极星的管理组
func GetManageGroupConfigForPolaris() ([]po.LcPerManagePermissionConfig, error) {
	configDataArr := make([]po.LcPerManagePermissionConfig, 0)
	configTypes := []interface{}{2, 3}
	dbErr := store.Mysql.SelectAllByCond(consts.TableManagePermissionConfig, db.Cond{
		consts.TcOrgId:   0,
		consts.TcType:    db.In(configTypes),
		consts.TcDelFlag: consts.AppIsNoDelete,
	}, &configDataArr)

	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}

	return configDataArr, nil
}
