package domain

import (
	"time"

	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/service/model/po"
	"upper.io/db.v3/lib/sqlbuilder"
)

// 初始化组织根部门
func InitOrgRootDept(orgId int64, operatorUid int64, orgName string, sourceChannel string, tx sqlbuilder.Tx) (int64, error) {
	logger.InfoF("[初始化根部门]-> orgId: %d,orgName: %s", orgId, orgName)
	deptBo := po.PpmOrgDepartment{
		Id:            snowflake.Id(),
		OrgId:         orgId,
		Name:          orgName,
		ParentId:      0,
		Sort:          0,
		SourceChannel: sourceChannel,
		Creator:       operatorUid,
		CreateTime:    time.Now(),
		Updator:       operatorUid,
		UpdateTime:    time.Now(),
	}
	dbErr := store.Mysql.TransInsert(tx, &deptBo)
	if dbErr != nil {
		logger.ErrorF("初始化根部门失败,原因: %s", dbErr)
		return 0, errs.MysqlOperateError
	}
	logger.InfoF("[初始化根部门]-> 成功  orgId: %d, deptId: %d", orgId, deptBo.Id)

	return deptBo.Id, nil
}

//初始化部门
func InitDepartment(orgId int64, outOrgId string, sourceChannel string, superAdminRoleId, normalAdminRoleId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	//if sourceChannel == consts.AppSourceChannelDingTalk {
	//	return InitDingTalkDepartment(orgId, outOrgId, superAdminRoleId, normalAdminRoleId, tx)
	//} else if sourceChannel == consts.AppSourceChannelFeiShu {
	//	return InitFsDepartment(orgId, outOrgId, superAdminRoleId, normalAdminRoleId, tx)
	//}
	return errs.BuildSystemErrorInfo(errs.SourceChannelNotDefinedError)
}

//func InitDingTalkDepartment(orgId int64, corpId string, superAdminRoleId, normalAdminRoleId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
//	deptList, err := dingtalk2.GetScopeDeps(corpId)
//	if err != nil {
//		return err
//	}
//	//把组织作为根部门插入
//	client, err := dingtalk2.GetDingTalkClientRest(corpId)
//	if err != nil {
//		logger.Error(err)
//		return errs.BuildSystemErrorInfo(errs.DingTalkClientError, err)
//	}
//
//	deptMap := map[int64]sdk.DepartmentInfo{}
//	for _, dep := range deptList {
//		deptMap[dep.Id] = dep
//	}
//
//	//判断根部门是否存在
//	rootDep, ok := deptMap[1]
//	if !ok {
//		depResp, rootErr := client.GetDeptDetail("1", nil)
//		if rootErr != nil {
//			logger.Error(rootErr)
//			return errs.BuildSystemErrorInfo(errs.DingTalkClientError, rootErr)
//		}
//		if depResp.ErrCode != 0 {
//			logger.Error(depResp.ErrMsg)
//			return errs.DingTalkClientError
//		}
//		rootDep = sdk.DepartmentInfo{
//			Id:              depResp.Id,
//			Name:            depResp.Name,
//			ParentId:        -1,
//			CreateDeptGroup: depResp.CreateDeptGroup,
//			AutoAddUser:     depResp.AutoAddUser,
//		}
//		deptList = append(deptList, rootDep)
//	}
//
//	depSize := len(deptList)
//	departmentInfo := make([]interface{}, len(deptList))
//	outDepartmentInfo := make([]interface{}, len(deptList))
//	fsDepIdMap := map[int64]int64{}
//
//	depIds := snowflake.IdBatch(depSize)
//
//	depOutIds := snowflake.IdBatch(depSize)
//
//	for k, v := range deptList {
//		fsDepIdMap[v.Id] = depIds[k]
//	}
//
//	rootId := fsDepIdMap[1]
//
//	for k, v := range deptList {
//		depId := depIds[k]
//		depOutId := depOutIds[k]
//
//		parentDepId := int64(0)
//
//		if id, ok := fsDepIdMap[v.ParentId]; ok {
//			parentDepId = id
//		} else {
//			parentDepId = rootId
//		}
//		if depId == rootId {
//			parentDepId = 0
//		}
//
//		departmentInfo[k] = &po.PpmOrgDepartment{
//			Id:            depId,
//			OrgId:         orgId,
//			Name:          v.Name,
//			ParentId:      parentDepId,
//			SourceChannel: consts.AppSourceChannelDingTalk,
//		}
//
//		outDepartmentInfo[k] = po.PpmOrgDepartmentOutInfo{
//			Id:                       depOutId,
//			OrgId:                    orgId,
//			DepartmentId:             depId,
//			SourceChannel:            consts.AppSourceChannelDingTalk,
//			OutOrgDepartmentId:       strconv.FormatInt(v.Id, 10),
//			Name:                     v.Name,
//			OutOrgDepartmentParentId: strconv.FormatInt(v.ParentId, 10),
//		}
//	}
//
//	//初始化用户
//	userInitErr := InitDingTalkUserList(orgId, corpId, fsDepIdMap, superAdminRoleId, normalAdminRoleId, tx)
//	if userInitErr != nil {
//		logger.Error(userInitErr)
//		return userInitErr
//	}
//
//	departErr := store.Mysql.TransBatchInsert(tx, &po.PpmOrgDepartment{}, departmentInfo)
//	if departErr != nil {
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, departErr)
//	}
//	outDepartErr := store.Mysql.TransBatchInsert(tx, &po.PpmOrgDepartmentOutInfo{}, outDepartmentInfo)
//	if outDepartErr != nil {
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, outDepartErr)
//	}
//	return nil
//}
//
//func InitFsDepartment(orgId int64, tenantKey string, superAdminRoleId, normalAdminRoleId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
//	deptList, err := feishu.GetScopeDeps(tenantKey)
//	if err != nil {
//		logger.Error(err)
//		return err
//	}
//	deptMap := map[string]vo.DepartmentRestInfoVo{}
//	for _, dep := range deptList {
//		deptMap[dep.Id] = dep
//	}
//
//	_, ok := deptMap["0"]
//	if !ok {
//		deptList = append(deptList, vo.DepartmentRestInfoVo{
//			Id:       "0",
//			Name:     "飞书平台组织",
//			ParentId: "Root_Department_Identification",
//		})
//	}
//
//	depSize := len(deptList)
//
//	departmentInfo := make([]interface{}, len(deptList))
//	outDepartmentInfo := make([]interface{}, len(deptList))
//	fsDepIdMap := map[string]int64{}
//
//	depIds := snowflake.IdBatch(depSize)
//
//	depOutIds := snowflake.IdBatch(depSize)
//
//	for k, v := range deptList {
//		fsDepIdMap[v.Id] = depIds[k]
//	}
//
//	rootId := fsDepIdMap["0"]
//
//	for k, v := range deptList {
//		depId := depIds[k]
//		depOutId := depOutIds[k]
//
//		parentDepId := int64(0)
//
//		if id, ok := fsDepIdMap[v.ParentId]; ok {
//			parentDepId = id
//		} else {
//			parentDepId = rootId
//		}
//		if depId == rootId {
//			parentDepId = 0
//		}
//
//		departmentInfo[k] = &po.PpmOrgDepartment{
//			Id:            depId,
//			OrgId:         orgId,
//			Name:          v.Name,
//			ParentId:      parentDepId,
//			SourceChannel: consts.AppSourceChannelFeiShu,
//		}
//
//		outDepartmentInfo[k] = po.PpmOrgDepartmentOutInfo{
//			Id:                       depOutId,
//			OrgId:                    orgId,
//			DepartmentId:             depId,
//			SourceChannel:            consts.AppSourceChannelFeiShu,
//			OutOrgDepartmentId:       v.Id,
//			Name:                     v.Name,
//			OutOrgDepartmentParentId: v.ParentId,
//		}
//	}
//
//	//初始化用户
//	userInitErr := InitFsUserList(orgId, tenantKey, fsDepIdMap, superAdminRoleId, normalAdminRoleId, tx)
//	if userInitErr != nil {
//		logger.Error(userInitErr)
//		return userInitErr
//	}
//
//	departErr := store.Mysql.TransBatchInsert(tx, &po.PpmOrgDepartment{}, departmentInfo)
//	if departErr != nil {
//		logger.Error(departErr)
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, departErr)
//	}
//	outDepartErr := store.Mysql.TransBatchInsert(tx, &po.PpmOrgDepartmentOutInfo{}, outDepartmentInfo)
//	if outDepartErr != nil {
//		logger.Error(outDepartErr)
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, outDepartErr)
//	}
//	return nil
//}
