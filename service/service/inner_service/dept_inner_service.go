package inner_service

import (
	"container/list"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/resp/inner_resp"
	"upper.io/db.v3"
)

/**
部门 内部调用
*/

// GetDeptList 根据Org获取部门列表 内部调用
func GetDeptList(orgId int64) ([]inner_resp.DeptInfoInnerResp, errs.SystemErrorInfo) {
	if orgId == 0 {
		return nil, errs.OrgNotExist
	}
	//去重
	deptBos, dbErr := domain.GetDeptListByOrg(orgId)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	deptList := make([]inner_resp.DeptInfoInnerResp, 0)
	for _, info := range deptBos {
		var dept inner_resp.DeptInfoInnerResp
		_ = copyer.Copy(info, &dept)
		deptList = append(deptList, dept)
	}
	return deptList, nil
}

// GetDeptListByIds 根据IDs获取部门列表 内部调用
func GetDeptListByIds(orgId int64, deptIds []int64, status, isDelete int) ([]inner_resp.DeptInfoInnerResp, errs.SystemErrorInfo) {
	if orgId == 0 {
		return nil, errs.OrgNotExist
	}
	if len(deptIds) == 0 {
		return []inner_resp.DeptInfoInnerResp{}, nil
	}
	//去重
	deptIds = slice.SliceUniqueInt64(deptIds)
	deptBos, dbErr := domain.GetDeptList(orgId, deptIds, nil, "", status, isDelete)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	deptList := make([]inner_resp.DeptInfoInnerResp, 0)
	for _, info := range deptBos {
		var dept inner_resp.DeptInfoInnerResp
		_ = copyer.Copy(info, &dept)
		deptList = append(deptList, dept)
	}
	return deptList, nil
}

func _copyBindBoToUserDeptData(infoBo *bo.UserDeptBindBo, infoResp *inner_resp.UserDeptBindData) {
	infoResp.DepartmentId = infoBo.DepartmentId
	infoResp.DepartmentName = infoBo.DepartmentName
	infoResp.OutOrgDepartmentId = infoBo.OutOrgDepartmentId
	infoResp.OutOrgDepartmentCode = infoBo.OutOrgDepartmentCode
	infoResp.PositionId = infoBo.OrgPositionId
	infoResp.PositionName = infoBo.PositionName
	infoResp.PositionLevel = infoBo.PositionLevel
	infoResp.IsLeader = infoBo.IsLeader
}

// GetUserDeptBindDataListByUsers 根据成员列表查询部门/职级列表
func GetUserDeptBindDataListByUsers(orgId int64, userIds []int64) (map[int64][]inner_resp.UserDeptBindData, errs.SystemErrorInfo) {
	// 查询部门/职级信息
	userDeptBindInfoList, dbErr := domain.GetUserDeptBindInfoListByUsers(orgId, userIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	respMap := make(map[int64][]inner_resp.UserDeptBindData)
	for _, infoBo := range userDeptBindInfoList {
		infoResp := inner_resp.UserDeptBindData{}
		_copyBindBoToUserDeptData(&infoBo, &infoResp)
		respMap[infoBo.UserId] = append(respMap[infoBo.UserId], infoResp)
	}
	return respMap, nil
}

func GetUserCountByDeptIds(orgID int64, deptIds []int64) (*inner_resp.GetUserCountByDeptIdsResp, errs.SystemErrorInfo) {
	// 获取部门树
	_, depMap, err := domain.GetDeptTree(orgID)
	if err != nil {
		logger.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.SystemError, err)
	}

	// 获取部门用户数统计
	depUserCounts, err := domain.GetAllDeptUserCount(orgID)
	if err != nil {
		logger.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.SystemError, err)
	}
	depUserCountMap := map[int64]uint64{}
	for _, depUserCount := range depUserCounts {
		depUserCountMap[depUserCount.DepartmentID] = depUserCount.Count
	}
	res := map[int64]uint64{}
	// 部门拼装统计参数
	for _, id := range deptIds {
		if _, ok := depUserCountMap[id]; !ok {
			res[id] = 0
			continue
		}
		depIds := map[int64]int8{}
		depIds[id] = 1

		depNode := depMap[id]
		if depNode == nil {
			res[id] = depUserCountMap[id]
			continue
		}

		depNode.Foreach(func(d *bo.DepartmentTreeNode) bool {
			if _, ok := depIds[d.ID]; ok {
				return false
			}
			depIds[d.ID] = 1
			return true
		})

		depUserCount := uint64(0)
		for depId, _ := range depIds {
			depUserCount = depUserCount + depUserCountMap[depId]
		}

		res[id] = depUserCount
	}

	return &inner_resp.GetUserCountByDeptIdsResp{UserCount: res}, nil
}

func GetUserDeptIds(orgId int64, userId int64) (*inner_resp.GetUserDeptIdsResp, errs.SystemErrorInfo) {
	info, err := domain.GetUserDeptIds(orgId, []int64{userId})
	if err != nil {
		logger.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.SystemError, err)
	}

	deptIds := []int64{}
	if depts, ok := info[userId]; ok {
		deptIds = depts
	}

	return &inner_resp.GetUserDeptIdsResp{DeptIds: deptIds}, nil
}

func GetUserDeptIdsBatch(orgId int64, userIds []int64) (*inner_resp.GetUserDeptIdsBatchResp, errs.SystemErrorInfo) {
	info, err := domain.GetUserDeptIds(orgId, userIds)
	if err != nil {
		logger.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.SystemError, err)
	}

	return &inner_resp.GetUserDeptIdsBatchResp{Data: info}, nil
}

func GetUserIdsByDeptIds(orgId int64, deptIds []int64) (*inner_resp.GetUserIdsByDeptIdsResp, errs.SystemErrorInfo) {
	// 获取部门树
	_, depMap, err := domain.GetDeptTree(orgId)
	if err != nil {
		logger.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ServerError, err)
	}

	allRelatedDeptIds := make([]int64, 0)
	depIds := map[int64]int8{}
	for _, deptId := range deptIds {
		allRelatedDeptIds = append(allRelatedDeptIds, deptId)
		deptNode := depMap[deptId]
		if deptNode != nil {
			deptNode.Foreach(func(d *bo.DepartmentTreeNode) bool {
				if _, ok := depIds[d.ID]; ok {
					return false
				}
				depIds[d.ID] = 1
				allRelatedDeptIds = append(allRelatedDeptIds, d.ID)
				return true
			})
		}
	}

	res := &inner_resp.GetUserIdsByDeptIdsResp{
		UserIds: []int64{},
	}
	if len(deptIds) == 0 {
		return res, nil
	}

	userIds, userIdsErr := domain.GetUserIdsByDeptIds(orgId, allRelatedDeptIds)
	if userIdsErr != nil {
		logger.Error(userIdsErr)
		return nil, userIdsErr
	}

	res.UserIds = userIds
	return res, nil
}

func GetLeadersByDeptIds(orgId int64, deptIds []int64) (*inner_resp.GetLeadersByDeptIdsResp, errs.SystemErrorInfo) {
	userDepts := make([]po.PpmOrgUserDepartment, 0)
	err := store.Mysql.SelectAllByCond(consts.TableUserDepartment, db.Cond{
		"department_id": db.In(deptIds),
		"is_leader":     1,
		"is_delete":     2,
		"org_id":        orgId,
	}, &userDepts)
	if err != nil {
		logger.Error(err)
		return nil, errs.MysqlOperateError
	}
	leaders := make([]inner_resp.DepartmentLeader, 0)
	for _, v := range userDepts {
		leaders = append(leaders, inner_resp.DepartmentLeader{
			DepartmentId: v.DepartmentId,
			LeaderId:     v.UserId,
		})
	}
	return &inner_resp.GetLeadersByDeptIdsResp{
		Leaders: leaders,
	}, nil
}

func GetDeptUserIds(orgId int64) (*inner_resp.GetDeptUserIdsResp, errs.SystemErrorInfo) {
	data, err := domain.GetDeptUserIdsMap(orgId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return &inner_resp.GetDeptUserIdsResp{Data: data}, nil
}

func GetFullDeptByIds(orgId int64, deptIds []int64) (map[int64][]string, errs.SystemErrorInfo) {
	resultMap := make(map[int64][]string, 0)
	if len(deptIds) < 1 {
		return resultMap, nil
	}
	// 先查询出该组织的所有部门树，再摘取出要查询的部门
	_, deMap, err := domain.GetDeptTree(orgId)
	if err != nil {
		logger.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	for _, deptId := range deptIds {
		depNode := deMap[deptId]
		resultMap[deptId] = GetFullDeptStrArrByTree(depNode)
	}

	return resultMap, nil
}

// GetFullDeptStrArrByTree 利用链表，将组装特定顺序的部门名称列表
func GetFullDeptStrArrByTree(deptNode *bo.DepartmentTreeNode) []string {
	strArr := make([]string, 0)
	link := list.New()

	if deptNode == nil {
		return strArr
	}
	cur := deptNode
	for cur != nil {
		link.PushFront(cur.Name)
		cur = cur.Parent
	}
	for element := link.Front(); element != nil; element = element.Next() {
		strArr = append(strArr, element.Value.(string))
	}

	return strArr
}

func GetUserDeptIdsWithParentId(orgId int64, userId int64) (*inner_resp.GetUserDeptIdsWithParentIdResp, errs.SystemErrorInfo) {
	depts, err := domain.GetUserDeptBindInfoListByUser(orgId, userId)
	if err != nil {
		logger.ErrorF("[GetUserDeptIdsWithParentId] GetUserDeptBindInfoListByUser err:%v, orgId:%v, userId:%v",
			err, orgId, userId)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	refDeptIds := make([]int64, 0, len(depts))
	for _, dept := range depts {
		refDeptIds = append(refDeptIds, dept.DepartmentId)
	}

	if len(refDeptIds) > 0 {
		parentIds, errSys := domain.GetDeptParentIds(orgId, refDeptIds)
		if errSys != nil {
			logger.ErrorF("[GetUserDeptIdsWithParentId] GetDeptParentIds err:%v, orgId:%v, uerId:%v",
				errSys, orgId, userId)
			return nil, errSys
		}
		refDeptIds = append(refDeptIds, parentIds...)
	}

	refDeptIds = slice.SliceUniqueInt64(refDeptIds)

	return &inner_resp.GetUserDeptIdsWithParentIdResp{DeptIds: refDeptIds}, nil
}
