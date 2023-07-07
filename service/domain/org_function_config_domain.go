package domain

import (
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/service/model/po"
	"upper.io/db.v3"
)

// GetFunctionConfig
func GetFunctionConfig(orgIds []int64) ([]po.PpmOrcFunctionConfig, error) {
	var pos []po.PpmOrcFunctionConfig
	dbErr := store.Mysql.SelectAllByCond(consts.TableOrgFunctionConfig, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    db.In(orgIds),
	}, &pos)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return pos, nil
}
