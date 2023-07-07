package service

import (
	"os"
	"strconv"

	"github.com/star-table/usercenter/core/conf"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
	"github.com/tealeg/xlsx"
)

/**
日志相关
*/

// GetLoginRecordsList
func GetLoginRecordsList(orgId int64, reqParam req.LoginRecordListReq) (*resp.LoginRecordListResp, errs.SystemErrorInfo) {
	count, records, dbErr := domain.GetLoginRecordList(orgId, reqParam)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	list := make([]resp.LoginRecordData, 0)
	for _, v := range records {
		accountNo := v.LoginName
		if accountNo == "" {
			accountNo = v.Mobile
		}
		if accountNo == "" {
			accountNo = v.Email
		}
		list = append(list, resp.LoginRecordData{
			ID:         v.Id,
			IP:         v.LoginIp,
			AccountNo:  accountNo,
			UserAgent:  v.UserAgent,
			Msg:        v.Msg,
			CreateTime: v.CreateTime,
		})
	}
	return &resp.LoginRecordListResp{
		List:  list,
		Total: count,
	}, nil
}

// ExportLoginRecordsList 导出登陆日志
func ExportLoginRecordsList(orgId, userId int64, reqParam req.ExportLoginRecordListReq) (string, errs.SystemErrorInfo) {
	info, err := GetLoginRecordsList(orgId, req.LoginRecordListReq{
		IP:        reqParam.IP,
		AccountNo: reqParam.AccountNo,
		StartDate: reqParam.StartDate,
		EndDate:   reqParam.EndDate,
		PageBo:    &bo.PageBo{Page: 1, Size: 100000},
	})
	if err != nil {
		logger.Error(err)
		return "", err
	}

	relatePath := "/records" + "/org_" + strconv.FormatInt(orgId, 10)
	excelDir := conf.Cfg.Resource.RootPath + relatePath
	mkdirErr := os.MkdirAll(excelDir, 0777)
	if mkdirErr != nil {
		logger.Error(mkdirErr)
		return "", errs.BuildSystemErrorInfo(errs.IssueDomainError, mkdirErr)
	}
	fileName := "登陆日志.xlsx"
	excelPath := excelDir + "/" + fileName
	url := conf.Cfg.Resource.LocalDomain + relatePath + "/" + fileName

	var file *xlsx.File
	var row *xlsx.Row
	var cell *xlsx.Cell

	file = xlsx.NewFile()
	sheet, ioErr := file.AddSheet("Sheet1")
	if ioErr != nil {
		logger.Error(ioErr)
		return "", errs.BuildSystemErrorInfo(errs.SystemError, ioErr)
	}

	row = sheet.AddRow()

	cell = row.AddCell()
	cell.Value = "流水编号"
	cell = row.AddCell()
	cell.Value = "后台账号"
	cell = row.AddCell()
	cell.Value = "操作时间"
	cell = row.AddCell()
	cell.Value = "操作内容"
	cell = row.AddCell()
	cell.Value = "IP地址"
	cell = row.AddCell()
	cell.Value = "登录设备"
	cell = row.AddCell()
	// 流水编号	后台账号	操作时间	操作内容	IP地址	登录设备
	for _, record := range info.List {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = strconv.FormatInt(record.ID, 10)
		cell = row.AddCell()
		cell.Value = record.AccountNo
		cell = row.AddCell()
		cell.Value = record.CreateTime.Format(consts.AppTimeFormat)
		cell = row.AddCell()
		cell.Value = record.Msg
		cell = row.AddCell()
		cell.Value = record.IP
		cell = row.AddCell()
		cell.Value = record.UserAgent
	}

	saveErr := file.Save(excelPath)
	if saveErr != nil {
		logger.Error(saveErr)
		return "", errs.SystemError
	}

	return url, nil
}
