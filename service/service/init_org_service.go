package service

import (
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InitOrg(initOrgBo bo.InitOrgBo, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	return domain.InitOrg(initOrgBo, tx)
}

func GeneralInitOrg(initOrgBo bo.InitOrgBo, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	return domain.GeneralInitOrg(initOrgBo, tx)
}
