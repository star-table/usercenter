package domain

import (
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"upper.io/db.v3"
)

func GetIndustryBoAllList(cond db.Cond) (*[]bo.IndustryBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmCmmIndustry{}
	err := store.Mysql.SelectAllByCond(consts.TableIndustry, cond, pos)
	if err != nil {
		logger.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.IndustryBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		logger.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, nil
}

func GetIndustryBo(industryId int64) (*bo.IndustryBo, errs.SystemErrorInfo) {
	info := &po.PpmCmmIndustry{}
	err := store.Mysql.SelectOneByCond(consts.TableIndustry, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       industryId,
	}, info)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.IndustryNotExist
		}
		return nil, errs.MysqlOperateError
	}

	res := &bo.IndustryBo{}
	_ = copyer.Copy(info, res)
	return res, nil
}
