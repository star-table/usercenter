package service

import (
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/resp"
	"upper.io/db.v3"
)

func IndustryList() (*resp.IndustryListResp, errs.SystemErrorInfo) {

	cond := db.Cond{
		consts.TcIsShow: consts.AppShowEnable,
	}

	bos, err := domain.GetIndustryBoAllList(cond)

	if err != nil {
		return nil, err
	}

	resultList := &[]*resp.IndustryResp{}

	copyErr := copyer.Copy(bos, resultList)
	if copyErr != nil {
		logger.ErrorF("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return &resp.IndustryListResp{
		List: *resultList,
	}, nil

}
