package open_service

import (
	"strconv"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/req/open_req"
	"github.com/star-table/usercenter/service/model/resp/open_resp"
	"upper.io/db.v3"
)

func _copyBindBoToUserDeptData(infoBo *bo.UserDeptBindBo, infoResp *open_resp.UserDeptBindData) {
	infoResp.DepartmentId = infoBo.DepartmentId
	infoResp.DepartmentName = infoBo.DepartmentName
	infoResp.OutOrgDepartmentId = infoBo.OutOrgDepartmentId
	infoResp.OutOrgDepartmentCode = infoBo.OutOrgDepartmentCode
	infoResp.PositionId = infoBo.OrgPositionId
	infoResp.PositionName = infoBo.PositionName
	infoResp.PositionLevel = infoBo.PositionLevel
	infoResp.IsLeader = infoBo.IsLeader
}

func _copyDeptBoToDeptInfoResp(infoBo *bo.OrgDeptBo, infoResp *open_resp.DeptInfoResp) {
	infoResp.Id = infoBo.Id
	infoResp.OrgId = infoBo.OrgId
	infoResp.Name = infoBo.Name
	infoResp.Code = infoBo.Code
	infoResp.ParentId = infoBo.ParentId
	infoResp.OutOrgDepartmentId = infoBo.OutOrgDepartmentId
	infoResp.OutOrgDepartmentCode = infoBo.OutOrgDepartmentCode
	infoResp.OutOrgDepartmentParentId = infoBo.OutOrgDepartmentParentId
	infoResp.Path = infoBo.Path
	infoResp.Sort = infoBo.Sort
	infoResp.IsHide = infoBo.IsHide
	infoResp.Status = infoBo.Status
}

// GetDeptList 根据Org获取部门列表
func GetDeptList(orgId int64) ([]open_resp.DeptInfoResp, errs.SystemErrorInfo) {
	deptBos, dbErr := domain.GetDeptListByOrg(orgId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	deptList := make([]open_resp.DeptInfoResp, 0)
	for _, info := range deptBos {
		dept := open_resp.DeptInfoResp{}
		_copyDeptBoToDeptInfoResp(&info, &dept)
		deptList = append(deptList, dept)
	}
	return deptList, nil
}

// GetDeptList 根据Org获取部门列表
func GetDeptListByQueryCond(orgId int64, reqParam open_req.DeptQueryReq) ([]open_resp.DeptInfoResp, errs.SystemErrorInfo) {
	deptBos, dbErr := domain.GetDeptListByLikeName(orgId, reqParam.LikeName)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	deptList := make([]open_resp.DeptInfoResp, 0)
	for _, info := range deptBos {
		dept := open_resp.DeptInfoResp{}
		_copyDeptBoToDeptInfoResp(&info, &dept)
		deptList = append(deptList, dept)
	}
	return deptList, nil
}

// GetDeptListByIds 根据IDs获取部门列表
func GetDeptListByIds(orgId int64, deptIds []int64) ([]open_resp.DeptInfoResp, errs.SystemErrorInfo) {
	//去重
	deptIds = slice.SliceUniqueInt64(deptIds)
	deptBos, dbErr := domain.GetDeptListByDeptIds(orgId, deptIds)

	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	deptList := make([]open_resp.DeptInfoResp, 0)
	for _, info := range deptBos {
		dept := open_resp.DeptInfoResp{}
		_copyDeptBoToDeptInfoResp(&info, &dept)
		deptList = append(deptList, dept)
	}
	return deptList, nil
}

// GetDeptChildrenList 获取子部门
func GetDeptChildrenList(orgId int64, deptIds []int64) ([]open_resp.DeptInfoResp, errs.SystemErrorInfo) {
	//去重
	childrenIds, dbErr := domain.GetDeptChildrenIdsByIds(orgId, deptIds)

	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	return GetDeptListByIds(orgId, childrenIds)
}

// GetUserDeptBindListByUser 根据成员id查询部门/职级列表
func GetUserDeptBindListByUser(orgId int64, userId int64) ([]open_resp.UserDeptBindResp, errs.SystemErrorInfo) {
	respList := make([]open_resp.UserDeptBindResp, 0)
	// 获取成员信息
	memberInfo, dbErr := domain.GetOrgMemberBaseInfoByUser(orgId, userId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return respList, nil
		}
		return nil, errs.MysqlOperateError
	}
	// 查询部门/职级信息
	userDeptBindInfoList, dbErr := domain.GetUserDeptBindInfoListByUser(orgId, userId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}

	for _, infoBo := range userDeptBindInfoList {
		infoResp := open_resp.UserDeptBindResp{}
		_copyBindBoToUserDeptData(&infoBo, &infoResp.UserDeptBindData)
		infoResp.UserId = memberInfo.UserId
		infoResp.Nickname = memberInfo.Name
		infoResp.Status = memberInfo.Status
		respList = append(respList, infoResp)
	}
	return respList, nil
}

// GetUserDeptBindListByUsers 根据成员列表查询部门/职级列表
func GetUserDeptBindListByUsers(orgId int64, userIds []int64) ([]open_resp.UserDeptBindResp, errs.SystemErrorInfo) {
	// 获取成员信息
	memberMap := make(map[int64]bo.OrgMemberBaseInfoBo)
	memberList, dbErr := domain.GetOrgMemberBaseInfoListByUsers(orgId, userIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	for _, member := range memberList {
		memberMap[member.UserId] = member
	}
	// 查询部门/职级信息
	userDeptBindInfoList, dbErr := domain.GetUserDeptBindInfoListByUsers(orgId, userIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	respList := make([]open_resp.UserDeptBindResp, 0)
	for _, infoBo := range userDeptBindInfoList {
		if member, ok := memberMap[infoBo.UserId]; ok {
			infoResp := open_resp.UserDeptBindResp{}
			_copyBindBoToUserDeptData(&infoBo, &infoResp.UserDeptBindData)
			infoResp.UserId = member.UserId
			infoResp.Nickname = member.Name
			infoResp.Status = member.Status
			respList = append(respList, infoResp)
		}
	}
	return respList, nil
}

// GetUserDeptBindDataListByUsers 根据成员列表查询部门/职级列表
func GetUserDeptBindDataListByUsers(orgId int64, userIds []int64) (map[int64][]open_resp.UserDeptBindData, errs.SystemErrorInfo) {
	// 查询部门/职级信息
	userDeptBindInfoList, dbErr := domain.GetUserDeptBindInfoListByUsers(orgId, userIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	respMap := make(map[int64][]open_resp.UserDeptBindData, 0)
	for _, infoBo := range userDeptBindInfoList {
		infoResp := open_resp.UserDeptBindData{}
		_copyBindBoToUserDeptData(&infoBo, &infoResp)
		respMap[infoBo.UserId] = append(respMap[infoBo.UserId], infoResp)
	}
	return respMap, nil
}

// GetSameOrSuperiorCount 获取平级或者上级的个数
func GetSameOrSuperiorCount(orgId int64, userId int64, superiorUid int64) (int, errs.SystemErrorInfo) {
	// 查询部门/职级信息
	userDeptBindInfoList, dbErr := domain.GetUserDeptBindInfoListByUsers(orgId, []int64{userId, superiorUid})
	if dbErr != nil {
		return 0, errs.MysqlOperateError
	}
	UserDeptBindMap := make(map[int64][]*bo.UserDeptBindBo)
	UserDeptBindUKMap := make(map[string]*bo.UserDeptBindBo)
	for _, infoBo := range userDeptBindInfoList {
		UserDeptBindMap[infoBo.DepartmentId] = append(UserDeptBindMap[infoBo.DepartmentId], &infoBo)
		key := strconv.FormatInt(infoBo.UserId, 10) + "_" + strconv.FormatInt(infoBo.DepartmentId, 10)
		UserDeptBindUKMap[key] = &infoBo
	}
	count := 0
	for deptId, infoBos := range UserDeptBindMap {
		// 小于2说明其中一个人不在该部门
		if len(infoBos) < 2 {
			continue
		}
		key1 := strconv.FormatInt(userId, 10) + "_" + strconv.FormatInt(deptId, 10)
		key2 := strconv.FormatInt(superiorUid, 10) + "_" + strconv.FormatInt(deptId, 10)
		u1Dept := UserDeptBindUKMap[key1]
		u2Dept := UserDeptBindUKMap[key2]
		// 平级或上级  越小等级越高
		if u2Dept.PositionLevel <= u1Dept.PositionLevel {
			count++
		}
	}
	return count, nil
}

// GetSameOrSuperiorCount 获取平级或者下级的成员Id
func GetSameOrSubordinateMembers(orgId int64, userId int64, deptIds []int64) ([]int64, errs.SystemErrorInfo) {
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	userIds := make([]int64, 0)
	// 先查询出来当前人员所在的部门
	userBindDeptList, dbErr := domain.GetUserDeptBindInfoListByUser(orgId, userId)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	// 循环部门
	for _, deptInfo := range userBindDeptList {
		// 最低级的，只获取自己
		if deptInfo.PositionLevel == consts.PositionMemberLevel {
			// 添加自己
			userIds = append(userIds, userId)
			break
		}

		// 所在部门及其子部门
		theAndChildren, dbErr := domain.GetDeptAndChildrenIds(orgId, []int64{deptInfo.DepartmentId})
		if dbErr != nil {
			logger.Error(dbErr)
			return nil, errs.MysqlOperateError
		}
		var userIdList []po.PpmUserId
		dbErr = conn.Select(
			"a.user_id",
		).
			From(db.Raw("(select "+
				" ud.user_id "+
				" from ppm_org_user_department ud "+
				" join ppm_org_department d "+
				" on ud.department_id = d.id "+
				" and d.org_id = ud.org_id "+
				" and d.id in ? "+
				" and d.is_delete = 2 "+
				" join lc_org_position p "+
				" on ud.org_position_id = p.org_position_id "+
				" and ud.org_id = p.org_id "+
				" and p.position_level >= ? "+
				" and p.is_delete = 2 "+
				" where ud.org_id = ? "+
				" and ud.is_delete = 2 ) a ", theAndChildren, deptInfo.PositionLevel, orgId)).All(&userIdList)
		if dbErr != nil {
			logger.Error(dbErr)
			return nil, errs.MysqlOperateError
		}
		for _, bind := range userIdList {
			userIds = append(userIds, bind.UserId)
		}

	}

	return slice.SliceUniqueInt64(userIds), nil
}
