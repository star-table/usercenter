package domain

import (
	"fmt"
	"strings"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp/inner_resp"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// GetDeptListByQuery 获取部门列表
func GetDeptListByQuery(orgId int64, params req.DepartmentListReq) ([]po.PpmOrgDepartment, int64, error) {
	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcStatus:   consts.AppStatusEnable,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcIsHide:   consts.AppIsNotHiding, //默认只查询非隐藏部门
	}
	var pos []po.PpmOrgDepartment
	//查询父部门的子部门信息
	if params.ParentID != nil {
		if *params.ParentID >= 0 {
			cond[consts.TcParentId] = params.ParentID
		}
	}
	//名称
	if params.Name != nil {
		cond[consts.TcName] = db.Like("%" + *params.Name + "%")
	}
	//查询顶级部门
	if params.IsTop != nil && *params.IsTop == 1 {
		cond[consts.TcParentId] = 0
	}
	//展示隐藏的部门
	if params.ShowHiding != nil && *params.ShowHiding == 1 {
		delete(cond, consts.TcIsHide)
	}
	if params.DeptIds != nil {
		if len(*params.DeptIds) > 0 {
			cond[consts.TcId] = db.In(*params.DeptIds)
		} else {
			return pos, 0, nil
		}
	}

	var page, size int64
	if params.Page != nil {
		page = *params.Page
	}
	if params.Size != nil {
		size = *params.Size
	}
	total, dbErr := store.Mysql.SelectAllByCondWithPageAndOrder(consts.TableDepartment, cond, nil, int(page), int(size), db.Raw("sort,create_time,id"), &pos)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, 0, dbErr
	}
	return pos, int64(total), nil
}

// UnBoundDepartmentUser 解绑部门用户，解绑当前用户所在的所有部门
func UnBoundDepartmentUser(orgId int64, userIds []int64, operatorId int64, tx sqlbuilder.Tx) error {
	//查询已有的绑定关系
	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   db.In(userIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	_, dbErr := store.Mysql.TransUpdateSmartWithCond(tx, consts.TableUserDepartment, cond, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  operatorId,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return dbErr
	}
	return nil
}

// GetDeptByDeptId 查找部门
func GetDeptByDeptId(orgId int64, deptId int64) (*bo.OrgDeptBo, error) {
	if deptId == 0 {
		return nil, db.ErrNoMoreRows
	}
	return _getDeptById(orgId, deptId, "")
}

// GetDeptByOutDeptId 查找部门
func GetDeptByOutDeptId(orgId int64, outDeptId string) (*bo.OrgDeptBo, error) {
	if outDeptId == "" {
		return nil, db.ErrNoMoreRows
	}
	return _getDeptById(orgId, 0, outDeptId)
}

// _getDeptById 查找部门
// deptId = 0 忽视条件
// outDeptId = "" 忽视条件
func _getDeptById(orgId int64, deptId int64, outDeptId string) (*bo.OrgDeptBo, error) {
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	cond := db.Cond{
		"d." + consts.TcOrgId:    orgId,
		"d." + consts.TcStatus:   consts.AppStatusEnable,
		"d." + consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if deptId != 0 {
		cond["d."+consts.TcId] = deptId
	}
	if outDeptId != "" {
		cond["od."+consts.TcOutOrgDepartmentId] = outDeptId
	}

	var dept bo.OrgDeptBo
	dbErr = conn.Select(
		"d.id",
		"d.org_id",
		"d.name",
		"d.code",
		"d.parent_id",
		"d.path",
		"d.sort",
		"d.is_hide",
		"d.source_channel",
		"d.source_platform",
		"d.status",
		"d.creator",
		"d.create_time",
		"d.updator",
		"d.update_time",
		"d.version",
		"d.is_delete",
		db.Raw("(case when `od`.`out_org_department_id` is null then '' else `od`.`out_org_department_id` end) as `out_org_department_id`"),
		db.Raw("(case when `od`.`out_org_department_code` is null then '' else `od`.`out_org_department_code` end) as `out_org_department_code`"),
		db.Raw("(case when `od`.`out_org_department_parent_id` is null then '' else `od`.`out_org_department_parent_id` end) as `out_org_department_parent_id`"),
	).From(consts.TableDepartment + " as d").
		LeftJoin(consts.TableDepartmentOutInfo + " as od").
		On(db.Raw(" d.id = od.department_id and d.org_id = od.org_id and od.is_delete = ?", consts.AppIsNoDelete)).
		Where(cond).Limit(1).One(&dept)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return &dept, nil
}

// GetDeptListByOrg 根据组织ID,获取部门列表
func GetDeptListByOrg(orgId int64) ([]bo.OrgDeptBo, error) {
	return _getDeptList(orgId, nil, nil, "", consts.AppStatusEnable, consts.AppIsNoDelete)
}

// GetDeptListByOrg 根据组织ID,获取部门列表
// likeName == "" 则忽视
func GetDeptListByLikeName(orgId int64, likeName string) ([]bo.OrgDeptBo, error) {
	return _getDeptList(orgId, nil, nil, likeName, consts.AppStatusEnable, consts.AppIsNoDelete)
}

// GetDeptListByDeptIds 根据部门ID列表,获取部门列表
// 注意 len(deptIds) == 0 则返回空切片
func GetDeptListByDeptIds(orgId int64, deptIds []int64) ([]bo.OrgDeptBo, error) {
	if len(deptIds) == 0 {
		return []bo.OrgDeptBo{}, nil
	}
	return _getDeptList(orgId, deptIds, nil, "", consts.AppStatusEnable, consts.AppIsNoDelete)
}

// GetDeptListByOutDeptIds 根据部门ID列表,获取部门列表
// 注意 len(outDeptIds) == 0 则返回空切片
func GetDeptListByOutDeptIds(orgId int64, outDeptIds []string) ([]bo.OrgDeptBo, error) {
	if len(outDeptIds) == 0 {
		return []bo.OrgDeptBo{}, nil
	}
	return _getDeptList(orgId, nil, outDeptIds, "", consts.AppStatusEnable, consts.AppIsNoDelete)
}

// GetDeptList 获取部门列表(所有条件)
func GetDeptList(orgId int64, deptIds []int64, outDeptIds []string, likeName string, status, isDelete int) ([]bo.OrgDeptBo, error) {
	if len(deptIds) == 0 {
		return []bo.OrgDeptBo{}, nil
	}
	return _getDeptList(orgId, deptIds, outDeptIds, likeName, status, isDelete)
}

// _getDeptList 获取部门列表
// len(deptIds) == 0 忽视条件
// len(outDeptIds) == 0 忽视条件
func _getDeptList(orgId int64, deptIds []int64, outDeptIds []string, likeName string, status, isDelete int) ([]bo.OrgDeptBo, error) {

	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	conn.SetLogging(true)
	cond := db.Cond{
		"d." + consts.TcOrgId: orgId,
	}
	if status > 0 {
		cond["d."+consts.TcStatus] = status
	}
	if isDelete > 0 {
		cond["d."+consts.TcIsDelete] = isDelete
	}
	if len(deptIds) > 0 {
		cond["d."+consts.TcId] = db.In(deptIds)
	}
	if len(outDeptIds) > 0 {
		cond["od."+consts.TcOutOrgDepartmentId] = db.In(outDeptIds)
	}
	if likeName != "" {
		cond["d."+consts.TcName] = db.Like("%" + likeName + "%")
	}
	var deptList []bo.OrgDeptBo
	dbErr = conn.Select(
		"d.id",
		"d.org_id",
		"d.name",
		"d.code",
		"d.parent_id",
		"d.path",
		"d.sort",
		"d.is_hide",
		"d.source_channel",
		"d.source_platform",
		"d.status",
		"d.creator",
		"d.create_time",
		"d.updator",
		"d.update_time",
		"d.version",
		"d.is_delete",
		db.Raw("(case when `od`.`out_org_department_id` is null then '' else `od`.`out_org_department_id` end) as `out_org_department_id`"),
		db.Raw("(case when `od`.`out_org_department_code` is null then '' else `od`.`out_org_department_code` end) as `out_org_department_code`"),
		db.Raw("(case when `od`.`out_org_department_parent_id` is null then '' else `od`.`out_org_department_parent_id` end) as `out_org_department_parent_id`"),
	).From(consts.TableDepartment+" as d").
		LeftJoin(consts.TableDepartmentOutInfo+" as od").
		On(db.Raw(" d.id = od.department_id and d.org_id = od.org_id and od.is_delete = ?", consts.AppIsNoDelete)).
		Where(cond).OrderBy("d.sort", "d.create_time", "d.id").All(&deptList)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return deptList, nil
}

// GetDeptChildrenIdsById 获取子部门Id
func GetDeptChildrenIdsById(orgId int64, deptId int64) ([]int64, error) {
	childrenIds, _, dbErr := _getDeptChildrenIds(orgId, []int64{deptId})
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return childrenIds, nil
}

// GetDeptChildrenIdsByIds 获取子部门Id
func GetDeptChildrenIdsByIds(orgId int64, deptIds []int64) ([]int64, error) {
	childrenIds, _, dbErr := _getDeptChildrenIds(orgId, deptIds)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	return childrenIds, nil
}

// GetDeptAndChildrenIds 查部门及其子部门Id
// 注意 len(deptIds) == 0 则返回空切片
func GetDeptAndChildrenIds(orgId int64, deptIds []int64) ([]int64, error) {
	if len(deptIds) == 0 {
		return []int64{}, nil
	}
	childrenIds, existIds, dbErr := _getDeptChildrenIds(orgId, deptIds)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	childrenIds = append(childrenIds, existIds...)
	return slice.SliceUniqueInt64(childrenIds), nil
}

// GetDeptAndChildrenBatch 获取一批部门的子部门 id
func GetDeptAndChildrenBatch(orgId int64, deptIds []int64) (map[int64][]int64, error) {
	resultMap := make(map[int64][]int64, 0)
	// 获取部门列表
	deptList, dbErr := GetDeptListByOrg(orgId)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	for _, deptId := range deptIds {
		childrenIds, existIds, err := _getDeptChildrenIdsWithAllDept(deptList, []int64{deptId})
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		childrenIds = append(childrenIds, existIds...)
		resultMap[deptId] = childrenIds
	}
	return resultMap, nil
}

// _getDeptChildrenIds 获取子部门Id
func _getDeptChildrenIds(orgId int64, deptIds []int64) ([]int64, []int64, error) {
	// 获取部门列表
	deptList, dbErr := GetDeptListByOrg(orgId)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, nil, dbErr
	}
	return _getDeptChildrenIdsWithAllDept(deptList, deptIds)
}

func _getDeptChildrenIdsWithAllDept(allDeptArr []bo.OrgDeptBo, deptIds []int64) ([]int64, []int64, error) {
	// 组装为map
	childrenMap := map[int64][]int64{}
	for _, dept := range allDeptArr {
		// 把当前节点加入父节点的子集合
		if dept.ParentId != 0 {
			children := childrenMap[dept.ParentId]
			if children == nil {
				children = make([]int64, 0)
			}
			children = append(children, dept.Id)
			childrenMap[dept.ParentId] = children
		}

		children := childrenMap[dept.Id]
		if children == nil {
			children = make([]int64, 0)
		}
		childrenMap[dept.Id] = children
	}
	childrenIds := make([]int64, 0)
	existIds := make([]int64, 0)
	for _, deptId := range deptIds {
		if _, ok := childrenMap[deptId]; ok {
			existIds = append(existIds, deptId)
			_loopChildrenIds(deptId, &childrenIds, childrenMap)
		}
	}
	return slice.SliceUniqueInt64(childrenIds), existIds, nil
}

// _getDeptParentIds 获取父部门Id
func GetDeptParentIds(orgId int64, deptIds []int64) ([]int64, errs.SystemErrorInfo) {
	// 获取部门列表
	deptList, dbErr := GetDeptListByOrg(orgId)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	return _getDeptParentIdsWithAllDept(deptList, deptIds), nil
}

func _getDeptParentIdsWithAllDept(allDeptArr []bo.OrgDeptBo, deptIds []int64) []int64 {
	// 组装为map
	parentIdMap := map[int64]int64{}
	for _, dept := range allDeptArr {
		parentIdMap[dept.Id] = dept.ParentId
	}
	deptIdMap := map[int64]bool{}
	for _, deptId := range deptIds {
		_loopParentIds(parentIdMap[deptId], deptIdMap, parentIdMap)
	}
	ids := make([]int64, 0)
	for k := range deptIdMap {
		ids = append(ids, k)
	}
	return ids
}

func _loopParentIds(parentId int64, deptIdMap map[int64]bool, parentIdMap map[int64]int64) {
	if parentId == 0 {
		return
	}
	if !deptIdMap[parentId] {
		deptIdMap[parentId] = true
		_loopParentIds(parentIdMap[parentId], deptIdMap, parentIdMap)
	}
}

// GetDeptUserCountBatch 获取部门下的人数统计,包含子部门
// 因为需要包含自部门,并且一个人可以关联到多个部门,无论这些部门是否是父子关系.因此,每个部门的人数需要分别进行 sql 查询.
// 谨慎:这个方法会有大量的 sql 查询
func GetDeptUserCountBatch(orgId int64, deptIdsMap map[int64][]int64) (map[int64]uint64, error) {
	countMap := make(map[int64]uint64, 0)
	// 获取部门列表
	for deptId, childrenIdsIncludeSelf := range deptIdsMap {
		tmpCount, err := GetDeptIdsUserCount(orgId, childrenIdsIncludeSelf)
		if err != nil {
			return countMap, err
		}
		countMap[deptId] = uint64(tmpCount)
	}
	return countMap, nil
}

func GetDeptIdsUserCount(orgId int64, deptIds []int64) (int64, error) {
	if deptIds == nil || len(deptIds) < 1 {
		return 0, nil
	}
	// 获取部门列表
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	userCountObj := map[string]int64{}
	dbErr = conn.Select(db.Raw("count(distinct(d.user_id)) user_count")).From(consts.TableUserDepartment+" as d", consts.TableUserOrganization+" as o").
		Where(db.Cond{
			"d." + consts.TcDepartmentId: db.In(deptIds),
			"d." + consts.TcIsDelete:     consts.AppIsNoDelete,
			"o." + consts.TcIsDelete:     consts.AppIsNoDelete,
			"o." + consts.TcOrgId:        orgId,
			"d." + consts.TcOrgId:        orgId,
			"d." + consts.TcUserId:       db.Raw("o." + consts.TcUserId),
			"o." + consts.TcStatus:       consts.AppStatusEnable,
		}).One(&userCountObj)
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, dbErr
	}
	return userCountObj["user_count"], nil
}

func _loopChildrenIds(parentId int64, childrenIds *[]int64, childrenMap map[int64][]int64) {
	if children, ok := childrenMap[parentId]; ok && len(children) > 0 {
		for _, id := range children {
			*childrenIds = append(*childrenIds, id)
			_loopChildrenIds(id, childrenIds, childrenMap)
		}
	}
}

// GetUserDeptBindInfoListByUser 获取成员绑定的部门/职级信息列表
func GetUserDeptBindInfoListByUser(orgId int64, userId int64) ([]bo.UserDeptBindBo, error) {
	if userId == 0 {
		return []bo.UserDeptBindBo{}, nil
	}
	return GetUserDeptBindInfoListByUsers(orgId, []int64{userId})
}

// GetUserDeptBindInfoListByUsers 获取成员绑定的部门/职级信息列表
func GetUserDeptBindInfoListByUsers(orgId int64, userIds []int64) ([]bo.UserDeptBindBo, error) {
	if len(userIds) == 0 {
		return []bo.UserDeptBindBo{}, nil
	}
	return _getUserDeptBindInfoList(orgId, nil, nil, userIds, 0)
}

// GetUserDeptBindInfoListByUsersRmDupli 获取成员绑定的部门/职级信息列表，去除多余重复的关联关系（去重）。
// 比如用户 a，有两条绑定部门 d 的数据，这时只返回一条
func GetUserDeptBindInfoListByUsersRmDupli(orgId int64, userIds []int64) ([]bo.UserDeptBindBo, error) {
	userDeptArr, err := GetUserDeptBindInfoListByUsers(orgId, userIds)
	if err != nil {
		return userDeptArr, err
	}
	newUserDeptArr := make([]bo.UserDeptBindBo, 0, len(userDeptArr))
	dupliMap1 := make(map[string]struct{}, len(userIds))
	for _, userDept := range userDeptArr {
		key := fmt.Sprintf("%v-%v", userDept.UserId, userDept.DepartmentId)
		if _, exist := dupliMap1[key]; exist {
			continue
		} else {
			dupliMap1[key] = struct{}{}
			newUserDeptArr = append(newUserDeptArr, userDept)
		}
	}

	return newUserDeptArr, nil
}

// GetUserDeptBindInfoListByDepts 获取成员绑定的部门/职级信息列表
func GetUserDeptBindInfoListByDept(orgId int64, deptId int64) ([]bo.UserDeptBindBo, error) {
	if deptId == 0 {
		return []bo.UserDeptBindBo{}, nil
	}
	return GetUserDeptBindInfoListByDepts(orgId, []int64{deptId})
}

// GetUserDeptBindInfoListByDepts 获取成员绑定的部门/职级信息列表
func GetUserDeptBindInfoListByDepts(orgId int64, deptIds []int64) ([]bo.UserDeptBindBo, error) {
	if len(deptIds) == 0 {
		return []bo.UserDeptBindBo{}, nil
	}
	return _getUserDeptBindInfoList(orgId, deptIds, nil, nil, 0)
}

// GetUserDeptBindInfoListByOutDept 获取成员绑定的部门/职级信息列表
func GetUserDeptBindInfoListByOutDept(orgId int64, outDeptId string) ([]bo.UserDeptBindBo, error) {
	if outDeptId == "" {
		return []bo.UserDeptBindBo{}, nil
	}
	return GetUserDeptBindInfoListByOutDepts(orgId, []string{outDeptId})
}

// GetUserDeptBindInfoListByOutDepts 获取成员绑定的部门/职级信息列表
func GetUserDeptBindInfoListByOutDepts(orgId int64, outDeptIds []string) ([]bo.UserDeptBindBo, error) {
	if len(outDeptIds) == 0 {
		return []bo.UserDeptBindBo{}, nil
	}
	return _getUserDeptBindInfoList(orgId, nil, outDeptIds, nil, 0)
}

// GetUserDeptLeaderBindInfoListByDept 获取成员绑定的部门/职级信息列表
func GetUserDeptLeaderBindInfoListByDept(orgId int64, deptId int64) ([]bo.UserDeptBindBo, error) {
	if deptId == 0 {
		return []bo.UserDeptBindBo{}, nil
	}
	return GetUserDeptLeaderBindInfoListByDepts(orgId, []int64{deptId})
}

// GetUserDeptLeaderBindInfoListByDepts 获取成员绑定的部门/职级信息列表
func GetUserDeptLeaderBindInfoListByDepts(orgId int64, deptIds []int64) ([]bo.UserDeptBindBo, error) {
	if len(deptIds) == 0 {
		return []bo.UserDeptBindBo{}, nil
	}
	return _getUserDeptBindInfoList(orgId, deptIds, nil, nil, consts.DepartmentIsLeader)
}

// GetUserDeptBindInfoListByOrg 获取成员绑定的部门/职级信息列表
func GetUserDeptBindInfoListByOrg(orgId int64) ([]bo.UserDeptBindBo, error) {
	return _getUserDeptBindInfoList(orgId, nil, nil, nil, 0)
}

// GetUserDeptBindInfoListByDeptsAndUserIds 根据多个部门和多个用户，获取这些用户在这些部门中的关联关系
func GetUserDeptBindInfoListByDeptsAndUserIds(orgId int64, deptIds, userIds []int64) ([]bo.UserDepartmentBo, error) {
	return _getUserDeptBindList(orgId, deptIds, userIds)
}

func _getUserDeptBindList(orgId int64, deptIds, userIds []int64) ([]bo.UserDepartmentBo, error) {
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	conn.SetLogging(true)
	cond := db.Cond{
		"ud." + consts.TcOrgId:    orgId,
		"ud." + consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if deptIds != nil && len(deptIds) > 0 {
		cond["ud."+consts.TcDepartmentId] = db.In(deptIds)
	}
	if userIds != nil && len(userIds) > 0 {
		cond["ud."+consts.TcUserId] = db.In(userIds)
	}
	var bindList []po.PpmOrgUserDepartment
	var bindBoList = make([]bo.UserDepartmentBo, 0)
	dbErr = conn.Select().From(consts.TableUserDepartment + " as ud").
		Where(cond).All(&bindList)
	copyer.Copy(bindList, &bindBoList)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}

	return bindBoList, nil
}

// _getUserDeptBindInfoList 获取部门/职级信息列表
func _getUserDeptBindInfoList(orgId int64, deptIds []int64, outDeptIds []string, userIds []int64, isLeader int) ([]bo.UserDeptBindBo, error) {
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	//conn.SetLogging(true)
	cond := db.Cond{
		"ud." + consts.TcOrgId:    orgId,
		"ud." + consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if len(deptIds) > 0 {
		cond["ud."+consts.TcDepartmentId] = db.In(deptIds)
	}
	if len(outDeptIds) > 0 {
		cond["od."+consts.TcOutOrgDepartmentId] = db.In(outDeptIds)
	}
	if len(userIds) > 0 {
		cond["ud."+consts.TcUserId] = db.In(userIds)
	}
	if consts.CheckLeaderValue(isLeader) {
		cond["ud."+consts.TcIsLeader] = isLeader
	}
	var bindList []bo.UserDeptBindBo
	var bindListValPtr []bo.UserDeptBindBoValPtr
	dbErr = conn.Select(
		"ud.org_id",
		"ud.user_id",
		"ud.department_id",
		"ud.is_leader",
		"d.name as department_name",
		"ud.org_position_id",
		"p.position_level",
		"p.name as position_name",
		db.Raw("(case when `od`.`out_org_department_id` is null then '' else `od`.`out_org_department_id` end) as `out_org_department_id`"),
		db.Raw("(case when `od`.`out_org_department_code` is null then '' else `od`.`out_org_department_code` end) as `out_org_department_code`"),
		db.Raw("(case when `od`.`out_org_department_parent_id` is null then '' else `od`.`out_org_department_parent_id` end) as `out_org_department_parent_id`"),
	).From(consts.TableUserDepartment + " as ud").
		Join(consts.TableDepartment + " as d").
		On(db.Raw(" ud.department_id = d.id and d.org_id = ud.org_id and d.is_delete = ?", consts.AppIsNoDelete)).
		// 此处之前是 Join，但使用 Join 会导致一些条件下查询不到数据，改成 LeftJoin 后，字段需要改成指针类型（所以有了 `bo.UserDeptBindBoValPtr`），便于 Scan（存储） Null 值。
		LeftJoin(consts.TablePosition + " as p").
		On(db.Raw(" ud.org_position_id = p.org_position_id and p.org_id = ud.org_id and p.is_delete = ?", consts.AppIsNoDelete)).
		LeftJoin(consts.TableDepartmentOutInfo + " as od").
		On(db.Raw(" d.id = od.department_id and d.org_id = od.org_id and od.is_delete = ?", consts.AppIsNoDelete)).
		Where(cond).All(&bindListValPtr)
	copyErr := copyer.Copy(bindListValPtr, &bindList)
	if copyErr != nil {
		logger.Error(copyErr)
		return nil, errs.ObjectCopyError
	}
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}

	return bindList, nil
}

// GetUserDeptBindInfoByUserAndDept 根据获取用户和部门绑定关系
func GetUserDeptBindInfoByUserAndDept(orgId int64, userId int64, deptId int64) (*bo.UserDeptBindBo, error) {
	conn, dbErr := store.Mysql.GetConnect()
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, dbErr
	}
	mid := conn.Select(
		"ud.org_id",
		"ud.user_id",
		"ud.department_id",
		"ud.is_leader",
		"d.name as department_name",
		"ud.org_position_id",
		"p.position_level",
		"p.name as position_name",
		db.Raw("(case when `od`.`out_org_department_id` is null then '' else `od`.`out_org_department_id` end) as `out_org_department_id`"),
		db.Raw("(case when `od`.`out_org_department_code` is null then '' else `od`.`out_org_department_code` end) as `out_org_department_code`"),
		db.Raw("(case when `od`.`out_org_department_parent_id` is null then '' else `od`.`out_org_department_parent_id` end) as `out_org_department_parent_id`"),
	).From(consts.TableUserDepartment + " as ud").
		Join(consts.TableDepartment + " as d").
		On(db.Raw(" ud.department_id = d.id and d.org_id = ud.org_id and d.is_delete = ?", consts.AppIsNoDelete)).
		Join(consts.TablePosition + " as p").
		On(db.Raw(" ud.org_position_id = p.org_position_id and p.org_id = ud.org_id and p.is_delete = ?", consts.AppIsNoDelete)).
		LeftJoin(consts.TableDepartmentOutInfo + " as od").
		On(db.Raw(" d.id = od.department_id and d.org_id = od.org_id and od.is_delete = ?", consts.AppIsNoDelete)).
		Where(db.Cond{
			"ud." + consts.TcOrgId:        orgId,
			"ud." + consts.TcUserId:       userId,
			"ud." + consts.TcDepartmentId: deptId,
			"ud." + consts.TcIsDelete:     consts.AppIsNoDelete,
		})
	var userDeptBindInfo bo.UserDeptBindBo
	dbErr = mid.Limit(1).One(&userDeptBindInfo)
	if dbErr != nil {
		return nil, dbErr
	}

	return &userDeptBindInfo, nil
}

// 获取部门用户ID数组，如果deptId为0表示查询根部门用户
// deptId > 0 查询指定部门用户
// deptId = 0 查询根部门用户
// deptId < 0 查询所有用户
func GetDeptUserIds(orgId, deptId int64, extraDeptIds []int64, params *bo.GetDeptUserIdsParams, page uint, size uint) ([]int64, bool, error) {
	if page == 0 {
		page = 1
	}
	hasMore := false
	userIds := make([]int64, 0)
	conn, err := store.Mysql.GetConnect()
	if err != nil {
		logger.Error(err)
		return nil, false, errs.MysqlOperateError
	}

	//conn.SetLogging(true)
	users := make([]po.PpmOrgUser, 0)
	selector := conn.Select(db.Raw("u.*")).
		From(db.Raw(consts.TableUserOrganization + " uo")).
		LeftJoin(db.Raw(consts.TableUser + " u")).
		On(db.Raw("u.id = uo.user_id")).
		Where(db.Cond{
			"uo." + consts.TcOrgId:       orgId,
			"uo." + consts.TcIsDelete:    consts.AppIsNoDelete,
			"uo." + consts.TcCheckStatus: consts.AppCheckStatusSuccess, //审核通过的
			"uo." + consts.TcStatus:      consts.AppStatusEnable,       //启用的用户
			"u." + consts.TcIsDelete:     consts.AppIsNoDelete,
		})

	extraCond := make([]db.Compound, 0)
	userIdScopeCond := make([]db.Compound, 0)
	if deptId == 0 {
		raw := db.Raw("select user_id from ppm_org_user_department where org_id = ? and is_delete = ?", orgId, consts.AppIsNoDelete)
		userIdScopeCond = append(userIdScopeCond, db.Cond{
			"u." + consts.TcId: db.NotIn(raw),
		})
	} else if deptId > 0 {
		raw := db.Raw("select user_id from ppm_org_user_department where org_id = ? and department_id = ? and is_delete = ?", orgId, deptId, consts.AppIsNoDelete)
		selector = selector.And(db.Cond{
			"u." + consts.TcId: db.In(raw),
		})
	}

	if len(extraDeptIds) > 0 {
		raw := db.Raw("select user_id from ppm_org_user_department where org_id = ? and is_delete = ? and department_id in ?", orgId, consts.AppIsNoDelete, extraDeptIds)
		extraCond = append(extraCond, db.Cond{
			" u." + consts.TcId: db.In(raw),
		})
	}
	if params != nil {
		if params.Query != nil && *params.Query != "" {
			namePy := strings.ToLower(*params.Query)
			selector = selector.And(db.Or(db.Cond{
				"u.name": db.Like("%" + *params.Query + "%"),
			}, db.Cond{
				db.Raw("lower(u." + consts.TcNamePinyin + ")"): db.Like("%" + namePy + "%"),
			}, db.Cond{
				db.Raw("u." + consts.TcMobile): db.Like("%" + *params.Query + "%"),
			}, db.Cond{
				db.Raw("u." + consts.TcEmail): db.Like("%" + *params.Query + "%"),
			}, db.Cond{
				db.Raw("u." + consts.TcLoginName): db.Like("%" + *params.Query + "%"),
			}))
		}
		// 如果存在用户id范围限制，如果范围为空数组，则认为没有用户匹配
		if params.UserIds != nil {
			if len(*params.UserIds) > 0 {
				userIdScopeCond = append(userIdScopeCond, db.Cond{
					"  u." + consts.TcId: db.In(*params.UserIds),
				})
			} else {
				userIdScopeCond = append(userIdScopeCond, db.Cond{
					"  u." + consts.TcId: -1,
				})
			}
		}
	}
	// 如果用户id条件不存在，则extraCond也没必要存在了
	if len(userIdScopeCond) > 0 {
		selector = selector.And(db.Or(db.And(userIdScopeCond...), db.And(extraCond...)))
	}

	selector = selector.OrderBy("u.id")
	paginator := selector.Paginate(size).Page(page)
	total, err := paginator.TotalEntries()
	if err != nil {
		logger.Error(err)
		return nil, false, errs.MysqlOperateError
	}
	if total <= uint64((page-1)*size) {
		return userIds, false, nil
	}
	if total > uint64(page*size) {
		hasMore = true
	}
	err = paginator.All(&users)
	if err != nil {
		logger.Error(err)
		return nil, false, errs.MysqlOperateError
	}
	for _, user := range users {
		userIds = append(userIds, user.Id)
	}
	return slice.SliceUniqueInt64(userIds), hasMore, nil
}

// 统计部门用户数量
func GetAllDeptUserCount(orgId int64) ([]bo.DepartUserCount, error) {
	conn, err := store.Mysql.GetConnect()
	if err != nil {
		logger.Error(err)
		return nil, errs.MysqlOperateError
	}
	depUserCounts := make([]bo.DepartUserCount, 0)
	err = conn.Select(db.Raw("count(distinct ud.user_id) as count, ud.department_id as department_id")).From(db.Raw("ppm_org_user_department ud, ppm_org_user_organization uo")).Where(db.Cond{
		"ud." + consts.TcOrgId:       orgId,
		"ud." + consts.TcIsDelete:    consts.AppIsNoDelete,
		"uo." + consts.TcIsDelete:    consts.AppIsNoDelete,
		"uo." + consts.TcCheckStatus: consts.AppCheckStatusSuccess, //审核通过的
		"uo." + consts.TcStatus:      consts.AppStatusEnable,       //启用的用户
		"ud." + consts.TcUserId:      db.Raw("uo." + consts.TcUserId),
	}).GroupBy("ud.department_id").All(&depUserCounts)
	if err != nil {
		logger.Error(err)
		return nil, errs.MysqlOperateError
	}
	rootDepUserCount := bo.DepartUserCount{}
	err = conn.Select(db.Raw("count(distinct user_id) as count")).From(consts.TableUserOrganization).Where(db.Cond{
		consts.TcOrgId:       orgId,
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcCheckStatus: consts.AppCheckStatusSuccess, //审核通过的
		consts.TcStatus:      consts.AppStatusEnable,       //启用的用户
	}).One(&rootDepUserCount)
	if err != nil && err != db.ErrNoMoreRows {
		logger.Error(err)
		return nil, errs.MysqlOperateError
	}
	return append(depUserCounts, rootDepUserCount), nil
}

// 获取组织下，所有部门对应的用户 id
func GetAllDeptUserId(orgId int64) (map[int64][]int64, error) {
	deps := make([]po.PpmOrgDepartment, 0)
	err := store.Mysql.SelectAllByCond(consts.TableDepartment, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &deps)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	deptIds := make([]int64, 0)
	for _, tmpDept := range deps {
		deptIds = append(deptIds, tmpDept.Id)
	}

	// 查询对应的用户 id 数据
	userDeptList := make([]po.PpmOrgUserDepartment, 0)
	err = store.Mysql.SelectAllByCondWithColumns(consts.TableUserDepartment, db.Raw("user_id, department_id"), db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcDepartmentId: db.In(deptIds),
		// 限定查询到的用户是有效的组织成员
		consts.TcUserId: db.In(db.Raw("select user_id from ppm_org_user_organization where is_delete in (0, 2) and check_status=? and `status`=? and org_id=?", consts.AppCheckStatusSuccess, consts.AppStatusEnable, orgId)),
	}, &userDeptList)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	deptUserMap := make(map[int64][]int64, 0)
	for _, item := range userDeptList {
		if len(deptUserMap[item.DepartmentId]) > 0 {
			deptUserMap[item.DepartmentId] = append(deptUserMap[item.DepartmentId], item.UserId)
		} else {
			deptUserMap[item.DepartmentId] = []int64{item.UserId}
		}
	}

	return deptUserMap, nil
}

// GetAllDeptIdsWithChildren 根据deptIds，找到所有部门Id以及子部门Id
func GetAllDeptIdsWithChildren(deptMap map[int64]*bo.DepartmentTreeNode, deptIds []int64) []int64 {
	idMap := make(map[int64]struct{})
	for _, id := range deptIds {
		dept := deptMap[id]
		dept.Walk(func(d *bo.DepartmentTreeNode) {
			idMap[d.ID] = struct{}{}
		})
	}
	var ids []int64
	for id, _ := range idMap {
		ids = append(ids, id)
	}
	return ids
}

// 获取部门树
func GetDeptTree(orgId int64) (*bo.DepartmentTreeNode, map[int64]*bo.DepartmentTreeNode, error) {
	deps := make([]po.PpmOrgDepartment, 0)
	err := store.Mysql.SelectAllByCond(consts.TableDepartment, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &deps)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}
	depMap := map[int64]*bo.DepartmentTreeNode{}
	for _, dep := range deps {
		depMap[dep.Id] = &bo.DepartmentTreeNode{
			ID:       dep.Id,
			Name:     dep.Name,
			ParentID: dep.ParentId,
		}
	}
	orgInfo, err := GetBaseOrgInfo("", orgId)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	depMap[0] = &bo.DepartmentTreeNode{
		Name:     orgInfo.OrgName,
		ParentID: -1,
		ID:       0,
	}
	for _, dep := range depMap {
		parent, ok := depMap[dep.ParentID]
		if ok {
			parent.Childs = append(parent.Childs, dep)
			dep.Parent = parent
		}
	}
	return depMap[0], depMap, nil
}

// 获取用户的部门关联
func GetUserDeptIds(orgId int64, userIds []int64) (map[int64][]int64, error) {
	userDeps := make([]po.PpmOrgUserDepartment, 0)
	err := store.Mysql.SelectAllByCond(consts.TableUserDepartment, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   db.In(userIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &userDeps)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	userDeptsMap := map[int64][]int64{}
	for _, userDept := range userDeps {
		userDeptsMap[userDept.UserId] = append(userDeptsMap[userDept.UserId], userDept.DepartmentId)
	}
	return userDeptsMap, nil
}

// 获取部门下的用户
func GetUserIdsByDeptIds(orgId int64, deptIds []int64) ([]int64, errs.SystemErrorInfo) {
	userList := make([]po.PpmOrgUserOrganization, 0)
	// 如果 deptIds 包含 0，表示根部门 —— 整个组织。
	// 因此，如果包含 0，可以直接查询整个组织的用户 id。
	containZero := false
	if exist, _ := slice.Contain(deptIds, int64(0)); exist {
		containZero = true
	}
	queryCond := db.Cond{
		consts.TcOrgId:       orgId,
		consts.TcCheckStatus: consts.AppCheckStatusSuccess,
		consts.TcStatus:      consts.AppStatusEnable,
		consts.TcIsDelete:    consts.AppIsNoDelete,
		// consts.TcUserId:      db.In(db.Raw("select user_id from ppm_org_user_department where org_id = ? and department_id in ? and is_delete = ?", orgId, deptIds, consts.AppIsNoDelete)),
	}
	if containZero {
		// no other cond
	} else {
		queryCond[consts.TcUserId] = db.In(db.Raw("select user_id from ppm_org_user_department where org_id = ? and department_id in ? and is_delete = ?", orgId, deptIds, consts.AppIsNoDelete))
	}
	err := store.Mysql.SelectAllByCondWithColumns(consts.TableUserOrganization, db.Raw("user_id"), queryCond, &userList)
	if err != nil {
		logger.Error(err)
		return nil, errs.MysqlOperateError
	}

	var userIds []int64
	for _, u := range userList {
		userIds = append(userIds, u.UserId)
	}

	return userIds, nil
}

func GetDeptSimpleInfo(orgId int64) ([]inner_resp.SimpleInfo, errs.SystemErrorInfo) {
	var deptInfo []po.PpmOrgDepartment
	err := store.Mysql.SelectAllByCond(consts.TableDepartment, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcStatus:   consts.AppStatusEnable,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcIsHide:   consts.AppIsNotHiding, //默认只查询非隐藏部门
	}, &deptInfo)
	if err != nil {
		logger.Error(err)
		return nil, errs.MysqlOperateError
	}
	res := make([]inner_resp.SimpleInfo, 0)
	for _, department := range deptInfo {
		res = append(res, inner_resp.SimpleInfo{
			Id:       department.Id,
			Name:     department.Name,
			ParentId: department.ParentId,
		})
	}
	return res, nil
}

func GetRepeatDeptInfo(orgId int64) ([]inner_resp.RepeatMemberInfo, errs.SystemErrorInfo) {
	var deptInfo []po.PpmOrgDepartment
	err := store.Mysql.SelectAllByCond(consts.TableDepartment, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcStatus:   consts.AppStatusEnable,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcIsHide:   consts.AppIsNotHiding, //默认只查询非隐藏部门
		consts.TcName:     db.In(db.Raw("select name from ppm_org_department where org_id = ? and status = ? and is_delete = 2 and is_hide = ? group by name having count(1) > 1", orgId, consts.AppStatusEnable, consts.AppIsNotHiding)),
	}, &deptInfo)
	if err != nil {
		logger.Error(err)
		return nil, errs.MysqlOperateError
	}

	if len(deptInfo) == 0 {
		return []inner_resp.RepeatMemberInfo{}, nil
	}

	parentIds := make([]int64, 0)
	deptMap := make(map[string][]po.PpmOrgDepartment, 0)
	for _, department := range deptInfo {
		if department.ParentId != int64(0) {
			parentIds = append(parentIds, department.ParentId)
		}
		deptMap[department.Name] = append(deptMap[department.Name], department)
	}
	parentIds = slice.SliceUniqueInt64(parentIds)
	parentMap := make(map[int64]string, 0)
	if len(parentIds) > 0 {
		var parentInfo []po.PpmOrgDepartment
		parentErr := store.Mysql.SelectAllByCond(consts.TableDepartment, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcStatus:   consts.AppStatusEnable,
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcIsHide:   consts.AppIsNotHiding, //默认只查询非隐藏部门
			consts.TcId:       db.In(parentIds),
		}, &parentInfo)
		if parentErr != nil {
			logger.Error(parentErr)
			return nil, errs.MysqlOperateError
		}

		for _, department := range parentInfo {
			parentMap[department.Id] = department.Name
		}
	}

	res := make([]inner_resp.RepeatMemberInfo, 0)
	for _, list := range deptMap {
		for _, department := range list {
			temp := inner_resp.RepeatMemberInfo{
				Id:         department.Id,
				Name:       department.Name,
				Department: []string{},
			}
			if parentName, ok := parentMap[department.ParentId]; ok {
				temp.Department = []string{parentName}
			}

			res = append(res, temp)
		}
	}

	return res, nil
}

func GetDeptUserIdsMap(orgId int64) (map[int64][]int64, errs.SystemErrorInfo) {
	var deptUsers []po.PpmOrgUserDepartment
	err := store.Mysql.SelectAllByCond(consts.TableUserDepartment, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
	}, &deptUsers)
	if err != nil {
		logger.Error(err)
		return nil, errs.MysqlOperateError
	}

	res := map[int64][]int64{}
	for _, user := range deptUsers {
		if ok, _ := slice.Contain(res[user.DepartmentId], user.UserId); !ok {
			res[user.DepartmentId] = append(res[user.DepartmentId], user.UserId)
		}
	}

	return res, nil
}

// DeptRemoveUsers 将多个用户移出某个部门
func DeptRemoveUsers(orgId int64, userIds []int64, deptId int64) (int64, error) {
	if userIds == nil || len(userIds) < 1 {
		return 0, nil
	}
	effectNum, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableUserDepartment, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcUserId:       db.In(userIds),
		consts.TcDepartmentId: deptId,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return effectNum, dbErr
	}
	return effectNum, nil
}

// DeleteExtraUserDeptRelations 检查一批用户的关联关系（用户-部门）是否有多余，如果多余则删除。
func DeleteExtraUserDeptRelations(orgId, deptId int64, userIds []int64, operatorUid int64, tx sqlbuilder.Tx) error {
	// group by 分组，检查是否存在多余的关联，需要修正
	userDeptReList := make([]po.PpmOrgUserDepartment, 0)
	if dbErr := store.Mysql.SelectAllByCond(consts.TableUserDepartment, db.Cond{
		consts.TcUserId:       db.In(userIds),
		consts.TcOrgId:        orgId,
		consts.TcDepartmentId: deptId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, &userDeptReList); dbErr != nil {
		logger.Error(dbErr)
		return dbErr
	}
	checkMultiMap := make(map[int64][]int64, 0)
	for _, item := range userDeptReList {
		if _, ok := checkMultiMap[item.UserId]; ok {
			checkMultiMap[item.UserId] = append(checkMultiMap[item.UserId], item.Id)
		} else {
			checkMultiMap[item.UserId] = []int64{item.Id}
		}
	}
	couldDeleteRowIds := make([]int64, 0)
	for _, dataIds := range checkMultiMap {
		// 超过 1，表示有多余的关联关系（用户-部门关联）
		if len(dataIds) > 1 {
			for i, rowId := range dataIds {
				if i != 0 {
					couldDeleteRowIds = append(couldDeleteRowIds, rowId)
				}
			}
		}
	}
	if len(couldDeleteRowIds) > 0 {
		_, dbErr := store.Mysql.TransUpdateSmartWithCond(tx, consts.TableUserDepartment, db.Cond{
			consts.TcId: db.In(couldDeleteRowIds),
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  operatorUid,
		})
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}
	}

	return nil
}
