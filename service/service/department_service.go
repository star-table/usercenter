package service

import (
	"sort"
	"strings"
	"time"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/snowflake"
	"github.com/star-table/usercenter/core/store"
	"github.com/star-table/usercenter/pkg/store/mysql"
	"github.com/star-table/usercenter/pkg/util"
	"github.com/star-table/usercenter/pkg/util/copyer"
	"github.com/star-table/usercenter/pkg/util/format"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/slice"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/po"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
	"github.com/star-table/usercenter/service/model/resp/inner_resp"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// DepartmentList 部门列表查询
func DepartmentList(orgId int64, params req.DepartmentListReq, perContext *inner_resp.OrgUserPerContext) (*resp.DepartmentList, errs.SystemErrorInfo) {

	departmentList, total, err := domain.GetDeptListByQuery(orgId, params)
	if err != nil {
		return nil, errs.MysqlOperateError
	}
	var resultList []*resp.Department
	_ = copyer.Copy(departmentList, &resultList)

	if len(resultList) == 0 {
		return &resp.DepartmentList{
			Total: total,
			List:  resultList,
		}, nil
	}
	// 检索部门
	contactDepList := make([]resp.ContactDepartment, 0)
	for _, dept := range departmentList {
		contactDepList = append(contactDepList, resp.ContactDepartment{
			ID:       dept.Id,
			Name:     dept.Name,
			ParentID: dept.ParentId,
			Visible:  dept.IsHide != 1,
		})
	}

	// 统计一下时长
	t1 := time.Now()
	// 获取部门树
	_, depMap, err := domain.GetDeptTree(orgId)
	if err != nil {
		logger.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	//获取所有的部门
	departmentIds := make([]int64, 0)
	for _, department := range resultList {
		departmentIds = append(departmentIds, department.ID)
	}
	// 查询所有部门的成员统计
	allDeptCountArr, err := domain.GetAllDeptUserCount(orgId)
	if err != nil {
		logger.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	allDeptCountMap := make(map[int64]uint64, 0)
	for _, item := range allDeptCountArr {
		allDeptCountMap[item.DepartmentID] = item.Count
	}
	// 查询所有部门对应用户 ids
	deptUserIdMap, err := domain.GetAllDeptUserId(orgId)
	if err != nil {
		logger.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	/*
		// 查询这些部门的子部门
		// 慎重：`domain.GetDeptAndChildrenBatch` 调用会在组织部门数很多时，耗时很长！
		childrenMap, err := domain.GetDeptAndChildrenBatch(orgId, departmentIds)
		if err != nil {
			logger.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		deptUserCountWithChildrenMap := make(map[int64]uint64, 0)
		// style 1
		for _, deptId := range departmentIds {
			sum := int(0)
			curDeptUserIds := make([]int64, 0)
			// 匹配出该部门下的所有子部门，统计其总人数
			if childrenIds, ok := childrenMap[deptId]; ok {
				for _, oneDeptId := range childrenIds {
					tmpUserIds := deptUserIdMap[oneDeptId]
					curDeptUserIds = append(curDeptUserIds, tmpUserIds...)
				}
				// 去重
				noRepeatedUserIds := slice.SliceUniqueInt64(curDeptUserIds)
				sum += len(noRepeatedUserIds)
			}
			deptUserCountWithChildrenMap[deptId] = uint64(sum)
		}
	*/

	// style 2
	//deptUserCountWithChildrenMap, dbErr := domain.GetDeptUserCountBatch(orgId, childrenMap)
	//if dbErr != nil {
	//	logger.Error(dbErr)
	//	return nil, errs.MysqlOperateError
	//}

	diff := time.Since(t1)
	logger.InfoF("部门统计人数,耗时记录 %.3f", diff.Seconds())

	//查找部门负责人
	leaders, dbErr := domain.GetUserDeptLeaderBindInfoListByDepts(orgId, departmentIds)
	if dbErr != nil {
		logger.Error(dbErr)
		return nil, errs.MysqlOperateError
	}

	userIds := make([]int64, 0)
	userDepartmentMap := map[int64][]int64{}
	for _, leader := range leaders {
		userIds = append(userIds, leader.UserId)
		userDepartmentMap[leader.DepartmentId] = append(userDepartmentMap[leader.DepartmentId], leader.UserId)
	}

	userIds = slice.SliceUniqueInt64(userIds)
	// 查询账号信息
	memberList, dbErr := domain.GetOrgMemberBaseInfoListByUsers(orgId, userIds)
	if dbErr != nil {
		return nil, errs.MysqlOperateError
	}
	userMap := map[int64]resp.DeptLeaderUserInfo{}
	for _, accountInfo := range memberList {
		userMap[accountInfo.UserId] = resp.DeptLeaderUserInfo{
			UserID: accountInfo.UserId,
			Name:   accountInfo.Name,
			NamePy: accountInfo.NamePinyin,
			Avatar: accountInfo.Avatar,
		}
	}

	for i, department := range resultList {
		// 查询部门及其子	部门的人数统计。
		deptIdMap := make(map[int64]int8, 0)
		deptIdMap[department.ID] = 1
		// 查找该部门下的子部门
		depNode := depMap[department.ID]
		if depNode == nil {
			resultList[i].DeptUserCount = int(allDeptCountMap[department.ID])
		} else {
			depNode.Foreach(func(d *bo.DepartmentTreeNode) bool {
				if _, ok := deptIdMap[d.ID]; ok {
					return false
				}
				deptIdMap[d.ID] = 1
				return true
			})
			// depUserCount := uint64(0)
			childDeptIds := make([]int64, 0)
			for depId, _ := range deptIdMap {
				// depUserCount = depUserCount + allDeptCountMap[depId]
				childDeptIds = append(childDeptIds, depId)
			}

			// resultList[i].DeptUserCount = int(allDeptCountMap[department.ID])
			resultList[i].DeptUserCount = GetDeptCountByDeptUserIds(deptUserIdMap, department.ID, childDeptIds)
		}

		resultList[i].LeaderInfo = []resp.DeptLeaderUserInfo{}
		// 是否可编辑
		if perContext.HasManageDept(department.ID) {
			resultList[i].Editable = true
			// 默认的根部门不可删除
			if department.ParentID != 0 {
				resultList[i].Deletable = true
			}

		}
		if _, ok := userDepartmentMap[department.ID]; ok {
			for _, i3 := range userDepartmentMap[department.ID] {
				if _, ok1 := userMap[i3]; ok1 {
					(resultList)[i].LeaderInfo = append(resultList[i].LeaderInfo, userMap[i3])
				}
			}
		}
	}

	return &resp.DepartmentList{
		Total: total,
		List:  resultList,
	}, nil
}

// GetDeptCountByDeptUserIds 通过部门用户id，获取该部门下的人数统计
func GetDeptCountByDeptUserIds(deptUserIdMap map[int64][]int64, targetDeptId int64, childDeptIds []int64) int {
	curDeptUserIds := make([]int64, 0)
	for _, oneDeptId := range childDeptIds {
		tmpUserIds := deptUserIdMap[oneDeptId]
		curDeptUserIds = append(curDeptUserIds, tmpUserIds...)
	}

	// 去重
	noRepeatedUserIds := slice.SliceUniqueInt64(curDeptUserIds)

	return len(noRepeatedUserIds)
}

// CreateDepartment 创建部门
func CreateDepartment(orgId int64, operatorUid int64, reqParam req.CreateDepartmentReq) (int64, errs.SystemErrorInfo) {
	logger.InfoF("[创建部门] -> reqParam: %s", json.ToJsonIgnoreError(reqParam))
	//检验部门名称
	if !format.VerifyDepartmentName(reqParam.Name) {
		return 0, errs.DepartmentNameInvalid
	}

	//查看部门名称是否重复
	isExist, dbErr := store.Mysql.IsExistByCond(consts.TableDepartment, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
		consts.TcName:     reqParam.Name,
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, errs.MysqlOperateError
	}
	if isExist {
		return 0, errs.DeptNameConflictErr
	}

	//查询父部门，如果父部门是 0，则创建顶级部门（默认顶级部门为 0）下的部门
	if reqParam.ParentID > 0 {
		_, dbErr = domain.GetDeptByDeptId(orgId, reqParam.ParentID)
		if dbErr != nil {
			if dbErr == db.ErrNoMoreRows {
				return 0, errs.ParentDeptNotExist
			}
			logger.Error(dbErr)
			return 0, errs.MysqlOperateError
		}
	}

	id := snowflake.Id()

	dbErr = store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		//新建部门
		res, dbErr := tx.InsertInto(consts.TableDepartment).
			Columns(consts.TcId,
				consts.TcOrgId,
				consts.TcName,
				consts.TcParentId,
				consts.TcSort,
				consts.TcCreator,
				consts.TcUpdator).
			Values(id,
				orgId,
				reqParam.Name,
				reqParam.ParentID,
				// 生成部门的sort
				db.Raw("(select case when max(sort) is null then 1 else max(sort) + 1 end from ppm_org_department a where a.org_id = ?)", orgId),
				operatorUid,
				operatorUid,
			).Exec()
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}
		_, dbErr = res.RowsAffected()
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}

		//添加部门负责人
		if len(reqParam.LeaderIds) > 0 {
			var departmentUserList []interface{}
			for _, leaderId := range reqParam.LeaderIds {
				newId := snowflake.Id()
				departmentUserList = append(departmentUserList, po.PpmOrgUserDepartment{
					Id:            newId,
					OrgId:         orgId,
					UserId:        leaderId,
					DepartmentId:  id,
					OrgPositionId: consts.PositionMemberId, // 默认成员
					IsLeader:      1,
					Creator:       operatorUid,
					Updator:       operatorUid,
				})
			}
			dbErr := store.Mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, departmentUserList)
			if dbErr != nil {
				logger.Error(dbErr)
				return dbErr
			}
		}

		return nil
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}

	return id, nil
}

// UpdateDepartment 修改部门信息
func UpdateDepartment(orgId, operatorUid int64, reqParam req.UpdateDepartmentReq) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[修改部门信息] -> orgId: %d,  reqParam: %s", orgId, json.ToJsonIgnoreError(reqParam))

	info, dbErr := domain.GetDeptByDeptId(orgId, reqParam.DepartmentId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.DepartmentNotExist
		}
		logger.Error(dbErr)
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	reqParam.Name = strings.TrimSpace(reqParam.Name)
	// 为空则不修改
	if reqParam.Name != "" && reqParam.Name != info.Name {
		//检验部门名称
		if !format.VerifyDepartmentName(reqParam.Name) {
			return false, errs.DepartmentNameInvalid
		}

		_, dbErr := store.Mysql.UpdateSmartWithCond(consts.TableDepartment, db.Cond{
			consts.TcId: reqParam.DepartmentId,
		}, mysql.Upd{
			consts.TcName:    reqParam.Name,
			consts.TcUpdator: operatorUid,
		})
		if dbErr != nil {
			logger.Error(dbErr)
			return false, errs.MysqlOperateError
		}
	}

	//查看之前的部门负责人
	var oldLeaders []po.PpmOrgUserDepartment
	dbErr = store.Mysql.SelectAllByCond(consts.TableUserDepartment, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcOrgId:        orgId,
		consts.TcDepartmentId: reqParam.DepartmentId,
	}, &oldLeaders)
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	var oldLeaderIds, oldUserIds []int64
	for _, department := range oldLeaders {
		if department.IsLeader == 1 {
			oldLeaderIds = append(oldLeaderIds, department.UserId)
		}
		oldUserIds = append(oldUserIds, department.UserId)
	}
	// 一个人不能既是主管，又是普通成员
	var interArr = make([]int64, 0)
	if reqParam.LeaderIds != nil && reqParam.UserIds != nil {
		interArr = slice.Int64Intersect(reqParam.LeaderIds, reqParam.UserIds)
	} else if reqParam.LeaderIds != nil && reqParam.UserIds == nil {
		interArr = slice.Int64Intersect(reqParam.LeaderIds, reqParam.UserIds)
	} else if reqParam.LeaderIds == nil && reqParam.UserIds != nil {
		interArr = slice.Int64Intersect(reqParam.LeaderIds, reqParam.UserIds)
	}
	if len(interArr) > 0 {
		return false, errs.UserIsLeaderAndUserDeny
	}

	dbErr = store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		// 更新部门主管
		if reqParam.LeaderIds != nil {
			if err := UpdateDeptLeaderWithMysql(orgId, operatorUid, reqParam.DepartmentId, oldLeaderIds, reqParam.LeaderIds, tx); err != nil {
				return err
			}
		}

		// 更新普通成员
		if reqParam.UserIds != nil {
			if err := UpdateDeptUser(orgId, operatorUid, reqParam.DepartmentId, oldUserIds, reqParam.UserIds, tx); err != nil {
				return err
			}
		}
		return nil
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	return true, nil
}

// UpdateDeptLeaderWithMysql 修改部门主管（增、减）
// 更新部门负责人时，可能是新增用户关联，也有可能是从普通部门成员**转变**为部门故责任
func UpdateDeptLeaderWithMysql(orgId, operatorUid int64, deptId int64, oldUserIds, newUserIds []int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	deleteIds, addIds := util.GetDifMemberIds(oldUserIds, newUserIds)
	if len(deleteIds) > 0 {
		// 检查是否存在多余的关联，如果有多余的，则清理掉。
		dbErr := domain.DeleteExtraUserDeptRelations(orgId, deptId, deleteIds, operatorUid, tx)
		if dbErr != nil {
			logger.Error(dbErr)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
		}
		_, dbErr = store.Mysql.TransUpdateSmartWithCond(tx, consts.TableUserDepartment, db.Cond{
			consts.TcUserId:       db.In(deleteIds),
			consts.TcOrgId:        orgId,
			consts.TcDepartmentId: deptId,
		}, mysql.Upd{
			consts.TcIsLeader: consts.DepartmentNotLeader,
			consts.TcUpdator:  operatorUid,
		})
		if dbErr != nil {
			logger.Error(dbErr)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
		}
	}

	if len(addIds) > 0 {
		// 查询是否已经是部门成员
		deptUserIds, _, dbErr := domain.GetDeptUserIds(orgId, deptId, nil, nil, 1, 20000)
		if dbErr != nil {
			logger.Error(dbErr)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
		}
		updUserIds := make([]int64, 0)
		addRelations := make([]interface{}, 0)
		for _, id := range addIds {
			if ok, _ := slice.Contain(oldUserIds, id); ok {
				//如果本来就是部门用户
				updUserIds = append(updUserIds, id)
			} else if ok2, _ := slice.Contain(deptUserIds, id); ok2 {
				updUserIds = append(updUserIds, id)
			} else {
				//如果没有要新增
				addRelations = append(addRelations, po.PpmOrgUserDepartment{
					Id:           snowflake.Id(),
					OrgId:        orgId,
					UserId:       id,
					DepartmentId: deptId,
					IsLeader:     1,
					Creator:      operatorUid,
				})
			}
		}
		if len(updUserIds) > 0 {
			_, dbErr := store.Mysql.TransUpdateSmartWithCond(tx, consts.TableUserDepartment, db.Cond{
				consts.TcUserId:       db.In(updUserIds),
				consts.TcOrgId:        orgId,
				consts.TcDepartmentId: deptId,
			}, mysql.Upd{
				consts.TcIsLeader: consts.DepartmentIsLeader,
				consts.TcUpdator:  operatorUid,
			})
			if dbErr != nil {
				logger.Error(dbErr)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
			}
		}
		if len(addRelations) > 0 {
			dbErr := store.Mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, addRelations)
			if dbErr != nil {
				logger.Error(dbErr)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
			}
		}
	}

	return nil
}

// UpdateDeptUser 修改部门普通成员（增、减）
func UpdateDeptUser(orgId, operatorUid int64, deptId int64, oldDeptUserIds, newDeptUserIds []int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	deleteIds, addIds := util.GetDifMemberIds(oldDeptUserIds, newDeptUserIds)
	if len(deleteIds) > 0 {
		// 删除部分普通成员
		_, dbErr := store.Mysql.TransUpdateSmartWithCond(tx, consts.TableUserDepartment, db.Cond{
			consts.TcUserId:       db.In(deleteIds),
			consts.TcOrgId:        orgId,
			consts.TcDepartmentId: deptId,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  operatorUid,
		})
		if dbErr != nil {
			logger.Error(dbErr)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
		}
	}
	if len(addIds) > 0 {
		// 新增部分普通成员。
		addRelations := make([]interface{}, 0)
		for _, id := range addIds {
			//如果没有要新增
			addRelations = append(addRelations, po.PpmOrgUserDepartment{
				Id:           snowflake.Id(),
				OrgId:        orgId,
				UserId:       id,
				DepartmentId: deptId,
				IsLeader:     consts.DepartmentNotLeader,
				Creator:      operatorUid,
				Updator:      operatorUid,
			})
		}
		if len(addRelations) > 0 {
			dbErr := store.Mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, addRelations)
			if dbErr != nil {
				logger.Error(dbErr)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
			}
		}
	}

	return nil
}

// DeleteDepartment 删除部门
func DeleteDepartment(orgId, operatorUid int64, reqParam req.DeleteDepartmentReq) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[删除部门] -> orgId: %d,  reqParam: %s", orgId, json.ToJsonIgnoreError(reqParam))
	_, dbErr := domain.GetDeptByDeptId(orgId, reqParam.DepartmentId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.DepartmentNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	//找部门下的子部门
	childrenIds, dbErr := domain.GetDeptChildrenIdsById(orgId, reqParam.DepartmentId)
	if dbErr != nil {
		return false, errs.MysqlOperateError
	}
	if len(childrenIds) > 0 {
		return false, errs.DeptHaveSubDeptError
	}
	// 为了兼容极星的虚拟的顶级部门缺失以及无码不允许删除顶级部门问题，这里先判断是否有多个顶级部门（parentId 为 0）
	// 如果有多个，则可以删除当前这个"顶级"部门，如果只有一个，则不能删除。
	topDeptIdArr, dbErr := domain.GetDeptChildrenIdsById(orgId, 0)
	if dbErr != nil {
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	if len(topDeptIdArr) == 1 {
		return false, errs.Top2DeptMustHaveOne
	}
	//if deptBo.ParentId == 0 {
	//	return false, errs.RootDeptCannotDeleteError
	//}
	dbErr = store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		//删除部门
		_, dbErr := store.Mysql.TransUpdateSmartWithCond(tx, consts.TableDepartment, db.Cond{
			consts.TcId:    reqParam.DepartmentId,
			consts.TcOrgId: orgId,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  operatorUid,
		})
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}
		//删除部门用户关系
		_, dbErr = store.Mysql.TransUpdateSmartWithCond(tx, consts.TableUserDepartment, db.Cond{
			consts.TcDepartmentId: reqParam.DepartmentId,
			consts.TcIsDelete:     consts.AppIsNoDelete,
			consts.TcOrgId:        orgId,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  operatorUid,
		})
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}
		// 删除管理组中绑定的部门关系
		_, dbErr = store.Mysql.TransUpdateSmartWithCond(tx, consts.TableManageGroup, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
			db.Raw("json_search(dept_ids, 'one', ?)", reqParam.DepartmentId): db.IsNotNull(),
		}, mysql.Upd{
			consts.TcDeptIds: db.Raw("json_remove(dept_ids,JSON_UNQUOTE(json_search(dept_ids, 'one', ?)))", reqParam.DepartmentId),
			consts.TcUpdator: operatorUid,
		})
		if dbErr != nil {
			logger.Error(dbErr)
			return dbErr
		}

		return nil
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	return true, nil
}

// ChangeUserDeptAndPosition 修改用户的部门/职级
func ChangeUserDeptAndPosition(orgId, operatorUid int64, reqParam req.ChangeUserDeptAndPositionReq, perContext *inner_resp.OrgUserPerContext) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[修改用户的部门/职级] -> orgId: %d,  reqParam: %s", orgId, json.ToJsonIgnoreError(reqParam))
	//查看用户是否存在
	_, dbErr := domain.GetEnableOrgMemberBaseInfoByUser(orgId, reqParam.UserId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.OrgMemberNotExistOrDisable
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	// 修改
	err := _changeUserDeptAndPosition(orgId, operatorUid, reqParam.UserId, reqParam.DeptAndPositions, perContext)
	if err != nil {
		return false, err
	}

	return true, nil
}

// AllocateUserDept 给用户分配部门
func AllocateUserDept(orgId, operatorUid int64, reqParam req.AllocateUserDeptReq, perContext *inner_resp.OrgUserPerContext) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[给用户分配部门] -> orgId: %d,  reqParam: %s", orgId, json.ToJsonIgnoreError(reqParam))
	//查看用户是否存在
	users, dbErr := domain.GetEnableOrgMemberBaseInfoListByUsers(orgId, reqParam.UserIds)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.OrgMemberNotExistOrDisable
		}
		logger.Error(dbErr)
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	userIds := make([]int64, 0)
	for _, item := range users {
		userIds = append(userIds, item.UserId)
	}
	userIds = slice.SliceUniqueInt64(userIds)
	//查看部门是否存在
	depts, dbErr := domain.GetDeptListByDeptIds(orgId, reqParam.DstDeptIds)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.DepartmentNotExist
		}
		logger.Error(dbErr)
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	deptIds := make([]int64, 0)
	for _, item := range depts {
		deptIds = append(deptIds, item.Id)
	}
	deptIds = slice.SliceUniqueInt64(deptIds)
	// 查询这批用户，在这批部门中，已经存在的关联关系
	userDeptBos, err := domain.GetUserDeptBindInfoListByDeptsAndUserIds(orgId, deptIds, userIds)
	if err != nil {
		logger.Error(err)
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	userDeptsMap := make(map[int64][]int64, 0)
	for _, item := range userDeptBos {
		if _, ok := userDeptsMap[item.UserId]; ok {
			userDeptsMap[item.UserId] = append(userDeptsMap[item.UserId], item.DepartmentId)
		} else {
			userDeptsMap[item.UserId] = []int64{item.DepartmentId}
		}
	}

	transErr := store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		//添加新的关联关系
		var departmentUser []interface{}
		for _, id := range userIds {
			// 当前用户已经所在的部门，已经存在该部门，则不需再添加关联关系。
			existDeptIds := make([]int64, 0)
			if tmpArr, ok := userDeptsMap[id]; ok {
				existDeptIds = tmpArr
			}
			for _, departmentId := range deptIds {
				if exist, _ := slice.Contain(existDeptIds, departmentId); !exist {
					departmentUser = append(departmentUser, po.PpmOrgUserDepartment{
						Id:           snowflake.Id(),
						OrgId:        orgId,
						UserId:       id,
						DepartmentId: departmentId,
						IsLeader:     consts.DepartmentNotLeader,
						Creator:      operatorUid,
						Updator:      operatorUid,
					})
				}
			}
		}
		if len(departmentUser) > 0 {
			insertErr := store.Mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, departmentUser)
			if insertErr != nil {
				logger.Error(insertErr)
				return insertErr
			}
		}

		return nil
	})
	if transErr != nil {
		logger.Error(transErr)
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, transErr)
	}

	return true, nil
}

// ChangeUserAdminGroup 切换用户管理组
// 目前只支持用户单个管理组的切换，如果一个用户有多个管理组，则此方法不适用
func ChangeUserAdminGroup(orgId, operatorUid int64, reqParam req.ChangeUserAdminGroupReq, perContext *inner_resp.OrgUserPerContext) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[切换用户管理组] -> orgId: %d,  reqParam: %s", orgId, json.ToJsonIgnoreError(reqParam))
	//查看用户是否存在
	_, dbErr := domain.GetEnableOrgMemberBaseInfoByUser(orgId, reqParam.UserId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.OrgMemberNotExistOrDisable
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	//查看管理组是否存在
	if reqParam.DstAdminGroupId > 0 {
		_, dbErr = domain.GetManageGroup(orgId, reqParam.DstAdminGroupId)
		if dbErr != nil {
			if dbErr == db.ErrNoMoreRows {
				return false, errs.ManageGroupNotExist
			}
			logger.Error(dbErr)
			return false, errs.MysqlOperateError
		}
	}

	dbErr = store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		oldAdminGroups, err := domain.GetManageGroupListByUsers(orgId, []int64{reqParam.UserId})
		if err != nil && err != db.ErrNoMoreRows {
			logger.Error(dbErr)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		// 检查用户是否是系统管理组，如果属于系统管理组，则不能更改他的角色。
		// 产品、测试：不能更改超管的管理组
		for _, oneGroup := range oldAdminGroups {
			if oneGroup.LangCode == consts.ManageGroupSys {
				// 产品：更换超管时，需要通过手机验证
				// 获取用户手机号
				dstUserInfo, err := domain.GetUserPoById(reqParam.UserId)
				if err != nil {
					logger.Error(err)
					return err
				}
				err = domain.AuthCodeVerify(consts.AuthCodeTypeChangeSuperAdmin, consts.ContactAddressTypeMobile, dstUserInfo.Mobile, reqParam.AuthCode)
				if err != nil {
					logger.Error(err)
					return err
				}
				// AuthCodeVerify 已经校验通过了

				tmpToken, err := domain.GetSMSLoginCode(consts.AuthCodeTypeChangeSuperAdmin, consts.ContactAddressTypeMobile, dstUserInfo.Mobile)
				if err != nil {
					logger.Error(err)
					return err
				}
				if tmpToken != reqParam.AuthCode {
					return errs.CaptchaError
				}
				// return errs.DenyChangeSysAdminGroupOfUser
			}
		}
		for _, oneGroup := range oldAdminGroups {
			tmpUserIds := make([]int64, 0)
			err = json.FromJson(oneGroup.UserIds, &tmpUserIds)
			if err != nil {
				logger.Error(err)
				return errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
			}
			// Search 参考 sort 包中的 Search 函数
			index := sort.Search(len(tmpUserIds), func(i int) bool { return tmpUserIds[i] == reqParam.UserId })
			// 删除旧管理组中特定的用户
			// eg: update lc_per_manage_group set user_ids = JSON_REMOVE(`user_ids`, '$[4]') WHERE org_id=1449 and id=1396673623212785665
			_, err = domain.RemoveUserFromAdminGroup(oneGroup.Id, operatorUid, index, tx)
			if err != nil {
				logger.Error(err)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}
		}
		if reqParam.DstAdminGroupId > 0 {
			_, err = domain.AppendUserIntoAdminGroup(reqParam.DstAdminGroupId, operatorUid, reqParam.UserId, tx)
			if err != nil {
				logger.Error(err)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}
		}
		return nil
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}

	return true, nil
}

// ChangeUserAdminGroup 切换用户管理组
func ChangeUserAdminGroupBatch(orgId, operatorUid int64, reqParam req.ChangeUserAdminGroupReq, perContext *inner_resp.OrgUserPerContext) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[切换用户管理组] -> orgId: %d,  reqParam: %s", orgId, json.ToJsonIgnoreError(reqParam))
	//查看用户是否存在
	_, dbErr := domain.GetEnableOrgMemberBaseInfoByUser(orgId, reqParam.UserId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.OrgMemberNotExistOrDisable
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	//查看管理组是否存在
	manageGroupList, manageGroupListErr := domain.GetManageGroupByCond(orgId, db.Cond{
		consts.TcId:       db.In(reqParam.DstAdminGroupIds),
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	})
	if manageGroupListErr != nil {
		logger.Error(manageGroupListErr)
		return false, errs.MysqlOperateError
	}
	if len(manageGroupList) != len(reqParam.DstAdminGroupIds) {
		return false, errs.ManageGroupNotExist
	}

	oldAdminGroups, err := domain.GetManageGroupListByUsers(orgId, []int64{reqParam.UserId})
	if err != nil && err != db.ErrNoMoreRows {
		logger.Error(dbErr)
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	oldGroupIds := []int64{}
	for _, group := range oldAdminGroups {
		oldGroupIds = append(oldGroupIds, group.Id)
	}
	deled, added := util.GetDifMemberIds(oldGroupIds, reqParam.DstAdminGroupIds)
	if len(deled) == 0 && len(added) == 0 {
		//没有变动
		return true, nil
	}
	//获取系统管理组
	sysGroup, sysGroupErr := domain.GetSysManageGroup(orgId)
	if sysGroupErr != nil {
		logger.Error(sysGroupErr)
		return false, errs.MysqlOperateError
	}

	if ok, _ := slice.Contain(deled, sysGroup.Id); ok {
		return false, errs.DenyChangeSysAdminGroupOfUser
	}
	if ok, _ := slice.Contain(added, sysGroup.Id); ok {
		return false, errs.DenyChangeSysAdminGroupOfUser
	}

	dbErr = store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		for _, oneGroup := range oldAdminGroups {
			if ok, _ := slice.Contain(deled, oneGroup.Id); !ok {
				continue
			}
			tmpUserIds := make([]int64, 0)
			err = json.FromJson(oneGroup.UserIds, &tmpUserIds)
			if err != nil {
				logger.Error(err)
				return errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
			}
			searchedIndex := domain.GetSearchedIndexArr(tmpUserIds, reqParam.UserId)
			if searchedIndex != -1 {
				// 删除旧管理组中特定的用户
				// eg: update lc_per_manage_group set user_ids = JSON_REMOVE(`user_ids`, '$[4]') WHERE org_id=1449 and id=1396673623212785665
				_, err = domain.RemoveUserFromAdminGroup(oneGroup.Id, operatorUid, searchedIndex, tx)
				if err != nil {
					logger.Error(err)
					return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
				}
			}

		}
		for _, i2 := range added {
			_, err = domain.AppendUserIntoAdminGroup(i2, operatorUid, reqParam.UserId, tx)
			if err != nil {
				logger.Error(err)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}
		}

		return nil
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}

	return true, nil
}

// _changeUserDeptAndPosition 修改用户的部门/职级 全量
func _changeUserDeptAndPosition(orgId, operatorUid int64, userId int64, deptAndPositions []req.DeptAndPositionReq, perContext *inner_resp.OrgUserPerContext) errs.SystemErrorInfo {

	// 部门ID
	deptIds := make([]int64, 0)
	deptPositionMap := make(map[int64]int64)

	for _, value := range deptAndPositions {
		if value.DepartmentId != 0 {
			deptIds = append(deptIds, value.DepartmentId)
			deptPositionMap[value.DepartmentId] = value.PositionId
		}
	}

	if len(deptIds) > 0 {
		//去重
		deptIds = slice.SliceUniqueInt64(deptIds)

		deptBos, dbErr := domain.GetDeptListByDeptIds(orgId, deptIds)
		if dbErr != nil {
			return errs.MysqlOperateError
		}
		// 个数不一致
		if len(deptIds) != len(deptBos) {
			return errs.DepartmentNotExist
		}
	}

	positionList, dbErr := domain.GetPositionList(orgId, consts.AppStatusEnable)
	if dbErr != nil {
		return errs.MysqlOperateError
	}
	positionMap := make(map[int64]po.LcOrgPosition)
	for _, p := range positionList {
		positionMap[p.OrgPositionId] = p
	}

	// 查询原关系
	oldBindList, dbErr := domain.GetUserDeptBindInfoListByUser(orgId, userId)
	if dbErr != nil {
		return errs.MysqlOperateError
	}
	// 按照部门分组map方便查询
	oldBindMap := make(map[int64]bo.UserDeptBindBo)
	for _, v := range oldBindList {
		oldBindMap[v.DepartmentId] = v
	}

	addBindList := make([]interface{}, 0)
	for _, deptId := range deptIds {
		// 非自己管理的
		if !perContext.HasManageDept(deptId) {
			continue
		}
		// 职级ID
		positionId := deptPositionMap[deptId]
		// 默认级别为成员
		if positionId == 0 {
			positionId = consts.PositionMemberId
		}
		// 添加原本不存在的, 对于修改了职级的直接删除原记录，重新添加
		if oldBindInfo, ok := oldBindMap[deptId]; !ok || positionId != oldBindInfo.OrgPositionId {
			userDept := po.PpmOrgUserDepartment{
				Id:            snowflake.Id(),
				OrgId:         orgId,
				DepartmentId:  deptId,
				UserId:        userId,
				OrgPositionId: positionId,
				Creator:       operatorUid,
				Updator:       operatorUid,
				IsDelete:      consts.AppIsNoDelete,
			}
			addBindList = append(addBindList, userDept)
		} else {
			// 存在的并且职级未修改的则无需操作
			delete(oldBindMap, deptId)
		}

	}

	// 把剩余的删除掉
	delBindDeptList := make([]int64, 0)
	for _, v := range oldBindMap {
		// 非自己管理的不进行删除
		if perContext.HasManageDept(v.DepartmentId) {
			delBindDeptList = append(delBindDeptList, v.DepartmentId)
		}
	}

	if len(addBindList) == 0 && len(delBindDeptList) == 0 {
		return nil
	}

	dbErr = store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		if len(delBindDeptList) > 0 {
			//删除旧有部门关联
			_, dbErr := store.Mysql.TransUpdateSmartWithCond(tx, consts.TableUserDepartment, db.Cond{
				consts.TcIsDelete:     consts.AppIsNoDelete,
				consts.TcDepartmentId: db.In(delBindDeptList),
				consts.TcUserId:       userId,
				consts.TcOrgId:        orgId,
			}, mysql.Upd{
				consts.TcIsDelete: consts.AppIsDeleted,
				consts.TcUpdator:  operatorUid,
			})
			if dbErr != nil {
				logger.Error(dbErr)
				return dbErr
			}
		}
		if len(addBindList) > 0 {
			//添加新的关联关系
			dbErr := store.Mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, addBindList)
			if dbErr != nil {
				logger.Error(dbErr)
				return dbErr
			}
		}

		return nil
	})

	if dbErr != nil {
		logger.Error(dbErr)
		return errs.MysqlOperateError
	}

	return nil
}

// UpdateUserDeptIds 更新用户的部门 todo
func UpdateUserDeptIds(orgId, operatorUid, userId int64, deptIds []int64, perContext *inner_resp.OrgUserPerContext) errs.SystemErrorInfo {
	if deptIds == nil {
		return nil
	}
	if len(deptIds) < 1 {
		// todo 清空用户关联的部门
	}
	//去重
	deptIds = slice.SliceUniqueInt64(deptIds)

	deptBos, dbErr := domain.GetDeptListByDeptIds(orgId, deptIds)
	if dbErr != nil {
		return errs.MysqlOperateError
	}
	// 个数不一致
	if len(deptIds) != len(deptBos) {
		return errs.DepartmentNotExist
	}
	// 查询已经存在的关联关系
	infos, dbErr := domain.GetUserDeptBindInfoListByDeptsAndUserIds(orgId, deptIds, []int64{userId})
	if dbErr != nil {
		return errs.MysqlOperateError
	}
	existDeptIds := make([]int64, 0)
	delDeptIds := make([]int64, 0)
	for _, item := range infos {
		existDeptIds = append(existDeptIds, item.DepartmentId)
	}
	addDeptIds := slice.ArrayDiffInt64(deptIds, existDeptIds)
	delDeptIds = slice.ArrayDiffInt64(existDeptIds, deptIds)
	if len(addDeptIds) < 1 {
		return nil
	}
	addBindList := make([]interface{}, 0)
	for _, deptId := range addDeptIds {
		userDept := po.PpmOrgUserDepartment{
			Id:            snowflake.Id(),
			OrgId:         orgId,
			DepartmentId:  deptId,
			UserId:        userId,
			OrgPositionId: 0,
			Creator:       operatorUid,
			Updator:       operatorUid,
			IsDelete:      consts.AppIsNoDelete,
		}
		addBindList = append(addBindList, userDept)
	}
	dbErr = store.Mysql.TransX(func(tx sqlbuilder.Tx) error {
		if len(addBindList) > 0 {
			//添加新的关联关系
			dbErr := store.Mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, addBindList)
			if dbErr != nil {
				logger.Error(dbErr)
				return dbErr
			}
		}
		if len(delDeptIds) > 0 {
			//删除旧有部门关联
			_, dbErr := store.Mysql.TransUpdateSmartWithCond(tx, consts.TableUserDepartment, db.Cond{
				consts.TcIsDelete:     consts.AppIsNoDelete,
				consts.TcDepartmentId: db.In(delDeptIds),
				consts.TcUserId:       userId,
				consts.TcOrgId:        orgId,
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
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	return nil
}

// UpdateDeptLeader 修改部门负责人
func UpdateDeptLeader(orgId, operatorUid int64, reqParam req.UpdateDeptLeaderReq) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[修改部门负责人] -> orgId: %d,  reqParam: %s", orgId, json.ToJsonIgnoreError(reqParam))

	if ok, _ := slice.Contain([]int{consts.DepartmentIsLeader, consts.DepartmentNotLeader}, reqParam.IsLeader); !ok {
		return false, errs.ParamError
	}
	//查看用户是否存在
	_, dbErr := domain.GetEnableOrgMemberBaseInfoByUser(orgId, reqParam.UserId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.OrgMemberNotExistOrDisable
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	// 查看部门是否存在
	_, dbErr = domain.GetDeptByDeptId(orgId, reqParam.DepartmentId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.DepartmentNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	// 查询部门用户是否存在
	bindInfo, dbErr := domain.GetUserDeptBindInfoByUserAndDept(orgId, reqParam.UserId, reqParam.DepartmentId)
	if dbErr != nil {
		// db错误
		if dbErr != db.ErrNoMoreRows {
			logger.Error(dbErr)
			return false, errs.MysqlOperateError
		}
		// 不存在则新增
		dbErr = store.Mysql.Insert(&po.PpmOrgUserDepartment{
			Id:            snowflake.Id(),
			OrgId:         orgId,
			UserId:        reqParam.UserId,
			DepartmentId:  reqParam.DepartmentId,
			IsLeader:      reqParam.IsLeader,
			OrgPositionId: consts.PositionMemberId,
			Creator:       operatorUid,
			Updator:       operatorUid,
		})
		if dbErr != nil {
			logger.Error(dbErr)
			return false, errs.MysqlOperateError
		}
	}

	if bindInfo.IsLeader == reqParam.IsLeader {
		return false, nil
	}
	// 存在则修改
	_, dbErr = store.Mysql.UpdateSmartWithCond(consts.TableUserDepartment, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcUserId:       bindInfo.UserId,
		consts.TcDepartmentId: bindInfo.DepartmentId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsLeader: reqParam.IsLeader,
		consts.TcUpdator:  operatorUid,
		consts.TcVersion:  db.Raw("version + 1"),
	})
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}
	return true, nil
}

// ChangeDeptSort 改变部门排序
func ChangeDeptSort(orgId, operatorUid int64, reqParam req.ChangeDeptSortReq) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[改变部门排序] -> orgId: %d,  reqParam: %s", orgId, json.ToJsonIgnoreError(reqParam))
	// 查询部门是否存在
	deptInfo, dbErr := domain.GetDeptByDeptId(orgId, reqParam.DeptId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.DepartmentNotExist
		}
		logger.Error(dbErr)
		return false, errs.MysqlOperateError
	}

	// 排序相同则无需改变
	if deptInfo.Sort == reqParam.Sort {
		return false, nil
	}

	// 改变
	_, dbErr = store.Mysql.UpdateSmartWithCond(consts.TableDepartment,
		db.Cond{
			consts.TcOrgId: orgId,
			consts.TcId:    deptInfo.Id,
		}, mysql.Upd{
			consts.TcSort:    reqParam.Sort,
			consts.TcUpdator: operatorUid,
			consts.TcVersion: db.Raw("version + 1"),
		})
	if dbErr != nil {
		logger.ErrorF("[改变部门排序] -> dbErr，原因:%s", dbErr)
		return false, errs.MysqlOperateError
	}

	return true, nil
}

// DeptRemoveUsers 将多个用户移出某个部门
func DeptRemoveUsers(orgId, operatorUid int64, reqParam req.DeptRemoveUsersReq) (bool, errs.SystemErrorInfo) {
	logger.InfoF("[将多个用户移出某个部门] -> orgId: %d,  reqParam: %s", orgId, json.ToJsonIgnoreError(reqParam))
	// 查询部门是否存在
	_, dbErr := domain.GetDeptByDeptId(orgId, reqParam.DeptId)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return false, errs.DepartmentNotExist
		}
		logger.Error(dbErr)
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	// 查询这些成员是否都在此部门中，如果有成员不在。
	effectNum, dbErr := domain.DeptRemoveUsers(orgId, reqParam.UserIds, reqParam.DeptId)
	if dbErr != nil {
		logger.Error(dbErr)
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	if effectNum < 1 {
		return false, errs.UsersNotInDept
	}

	return true, nil
}
