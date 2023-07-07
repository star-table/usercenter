package service

import (
	"strings"

	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/bo"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
)

func ContactFilter(orgID int64, reqObject req.ContactFilterReq) (*resp.ContactFilterResp, error) {
	departmentId := int64(0)
	if reqObject.ParentID != nil {
		departmentId = *reqObject.ParentID
	}

	params := req.DepartmentListReq{
		ParentID: &departmentId,
		DeptIds:  reqObject.Scope.DepartmentIds,
	}
	allDepts, _, err := domain.GetDeptListByQuery(orgID, params)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// 检索部门
	contactDepList := make([]resp.ContactDepartment, 0)
	for _, dept := range allDepts {
		contactDepList = append(contactDepList, resp.ContactDepartment{
			ID:       dept.Id,
			Name:     dept.Name,
			ParentID: dept.ParentId,
			Visible:  dept.IsHide != 1,
		})
	}

	// 获取部门树
	depTree, depMap, err := domain.GetDeptTree(orgID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// 检索当前部门下的用户
	contactUserList := make([]resp.ContactUser, 0)
	if reqObject.MemberLimit == 0 {
		reqObject.MemberLimit = 1000
	}
	if reqObject.MemberLimit > 1000 {
		reqObject.MemberLimit = 1000
	}
	page := reqObject.MemberOffset/reqObject.MemberLimit + 1
	size := reqObject.MemberLimit

	// 获取需要查询的额外部门用户
	extraDepIds := make([]int64, 0)
	if reqObject.OnlyMember {
		depIds := map[int64]int8{}
		for _, dept := range allDepts {
			extraDepIds = append(extraDepIds, dept.Id)
			deptNode := depMap[dept.Id]
			if deptNode != nil {
				deptNode.Foreach(func(d *bo.DepartmentTreeNode) bool {
					if _, ok := depIds[d.ID]; ok {
						return false
					}
					depIds[d.ID] = 1
					extraDepIds = append(extraDepIds, d.ID)
					return true
				})
			}
		}
	}

	getDeptUserIdsParams := &bo.GetDeptUserIdsParams{
		UserIds: reqObject.Scope.UserIds,
	}
	userIds, hasMore, err := domain.GetDeptUserIds(orgID, departmentId, extraDepIds, getDeptUserIdsParams, page, size)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	userPos, err := domain.GetUserListByIds(userIds)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	for _, userPo := range userPos {
		contactUserList = append(contactUserList, resp.ContactUser{
			ID:     userPo.Id,
			Name:   userPo.Name,
			Avatar: userPo.Avatar,
			Mobile: userPo.Mobile,
			Email:  userPo.Email,
		})
	}

	// 获取部门用户数统计
	depUserCounts, err := domain.GetAllDeptUserCount(orgID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	depUserCountMap := map[int64]uint64{}
	for _, depUserCount := range depUserCounts {
		depUserCountMap[depUserCount.DepartmentID] = depUserCount.Count
	}

	// 部门拼装统计参数
	userTotalNum := uint64(0)
	userTotalNum = userTotalNum + depUserCountMap[departmentId]
	for i, contactDep := range contactDepList {
		depIds := map[int64]int8{}
		depIds[contactDep.ID] = 1

		depNode := depMap[contactDep.ID]
		if depNode == nil {
			contactDepList[i].UserCount = uint32(depUserCountMap[contactDep.ID])
			continue
		}

		pathNodes, err := getDepPathNodes(depMap, contactDep.ID, false)
		if err != nil {
			logger.Error(err)
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
		contactDepList[i].UserCount = uint32(depUserCount)
		contactDepList[i].PathNodes = pathNodes
		contactDepList[i].ChildrenCount = uint32(len(depNode.Childs))
		userTotalNum = userTotalNum + depUserCount
	}

	// 拼装用户参数
	userDeptIds, err := domain.GetUserDeptIds(orgID, userIds)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	deptPathNodeMap := map[int64][]resp.ContactPathNode{}
	for i, contactUser := range contactUserList {
		if deptIds, ok := userDeptIds[contactUser.ID]; ok && len(deptIds) > 0 {
			pathNodesList := make([][]resp.ContactPathNode, 0)
			for _, deptId := range deptIds {
				if pathNodes, exist := deptPathNodeMap[deptId]; exist {
					pathNodesList = append(pathNodesList, pathNodes)
				} else {
					pathNodes, err := getDepPathNodes(depMap, deptId, true)
					if err != nil {
						logger.Error(err)
						continue
					}
					pathNodesList = append(pathNodesList, pathNodes)
					deptPathNodeMap[deptId] = pathNodes
				}
			}
			contactUserList[i].PathNodesList = pathNodesList
		}
	}

	//pathNodes
	pathNodes := make([]resp.ContactPathNode, 0)
	if departmentId <= 0 {
		pathNodes = append(pathNodes, resp.ContactPathNode{
			DepID:    depTree.ID,
			DepName:  depTree.Name,
			ParentID: depTree.ParentID,
		})
	} else {
		pathNodes, err = getDepPathNodes(depMap, departmentId, true)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	//获取组织总人数
	orgUserCount, err := domain.GetOrgUserCount(orgID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if reqObject.OnlyMember {
		contactDepList = make([]resp.ContactDepartment, 0)
	}

	return &resp.ContactFilterResp{
		DepList: contactDepList,
		// 部门层级，顺序返回
		PathNodes: pathNodes,
		// 是否有更多的用户
		UserHasMore: hasMore,
		// 用户总数
		UserTotalNum: userTotalNum,
		// 组织总人数
		OrgTotalNum: orgUserCount,
		// 用户列表
		User: contactUserList,
	}, nil
}

func getDepPathNodes(depMap map[int64]*bo.DepartmentTreeNode, deptId int64, containSelf bool) ([]resp.ContactPathNode, error) {
	depNode := depMap[deptId]
	if depNode == nil {
		return nil, nil
	}
	pathNodes := make([]resp.ContactPathNode, 0)
	if containSelf {
		pathNodes = append(pathNodes, resp.ContactPathNode{
			DepID:    depNode.ID,
			DepName:  depNode.Name,
			ParentID: depNode.ParentID,
			Visible:  true,
		})
	}
	parent := depNode.Parent
	loopMap := map[int64]int8{}
	for {
		if parent == nil {
			break
		}
		if _, ok := loopMap[parent.ID]; ok {
			break
		}
		pathNodes = append(pathNodes, resp.ContactPathNode{
			DepID:    parent.ID,
			DepName:  parent.Name,
			ParentID: parent.ParentID,
			Visible:  true,
		})
		loopMap[parent.ID] = 1
		parent = parent.Parent
	}
	for i, j := 0, len(pathNodes)-1; i < j; i, j = i+1, j-1 {
		pathNodes[i], pathNodes[j] = pathNodes[j], pathNodes[i]
	}
	return pathNodes, nil
}

func ContactSearch(orgID int64, searchReq req.ContactSearchReq) (*resp.ContactSearchResp, error) {
	query := strings.TrimSpace(searchReq.Query)

	if searchReq.SearchType == 2 {
		params := req.DepartmentListReq{
			Name:    &query,
			DeptIds: searchReq.Scope.DepartmentIds,
		}
		allDepts, _, err := domain.GetDeptListByQuery(orgID, params)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		if len(allDepts) == 0 {
			return &resp.ContactSearchResp{}, nil
		}

		// 检索部门
		contactDepList := make([]resp.ContactDepartment, 0)
		for _, dept := range allDepts {
			contactDepList = append(contactDepList, resp.ContactDepartment{
				ID:       dept.Id,
				Name:     dept.Name,
				ParentID: dept.ParentId,
				Visible:  dept.IsHide != 1,
			})
		}

		// 获取部门树
		_, depMap, err := domain.GetDeptTree(orgID)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		// 获取部门用户数统计
		depUserCounts, err := domain.GetAllDeptUserCount(orgID)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		depUserCountMap := map[int64]uint64{}
		for _, depUserCount := range depUserCounts {
			depUserCountMap[depUserCount.DepartmentID] = depUserCount.Count
		}

		// 部门拼装统计参数
		for i, contactDep := range contactDepList {
			depIds := map[int64]int8{}
			depIds[contactDep.ID] = 1

			depNode := depMap[contactDep.ID]
			if depNode == nil {
				contactDepList[i].UserCount = uint32(depUserCountMap[contactDep.ID])
				continue
			}

			pathNodes, err := getDepPathNodes(depMap, contactDep.ID, false)
			if err != nil {
				logger.Error(err)
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
			contactDepList[i].UserCount = uint32(depUserCount)
			contactDepList[i].PathNodes = pathNodes
			contactDepList[i].ChildrenCount = uint32(len(depNode.Childs))
		}
		return &resp.ContactSearchResp{
			DepList: contactDepList,
		}, nil
	} else if searchReq.SearchType == 1 {
		contactUserList := make([]resp.ContactUser, 0)
		if searchReq.Limit == 0 {
			searchReq.Limit = 1000
		}
		if searchReq.Limit > 1000 {
			searchReq.Limit = 1000
		}
		page := searchReq.Offset/searchReq.Limit + 1
		size := searchReq.Limit

		// 获取部门树
		_, depMap, err := domain.GetDeptTree(orgID)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		// 获取需要查询的额外部门用户
		extraDepIds := make([]int64, 0)
		if searchReq.OnlyMember {
			depIds := map[int64]int8{}
			if searchReq.Scope.DepartmentIds != nil {
				for _, deptId := range *searchReq.Scope.DepartmentIds {
					extraDepIds = append(extraDepIds, deptId)
					deptNode := depMap[deptId]
					if deptNode != nil {
						deptNode.Foreach(func(d *bo.DepartmentTreeNode) bool {
							if _, ok := depIds[d.ID]; ok {
								return false
							}
							depIds[d.ID] = 1
							extraDepIds = append(extraDepIds, d.ID)
							return true
						})
					}
				}
			}
		}

		userIds, hasMore, err := domain.GetDeptUserIds(orgID, -1, extraDepIds, &bo.GetDeptUserIdsParams{
			Query:   &query,
			UserIds: searchReq.Scope.UserIds,
		}, page, size)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		userPos, err := domain.GetUserListByIds(userIds)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		for _, userPo := range userPos {
			contactUserList = append(contactUserList, resp.ContactUser{
				ID:        userPo.Id,
				Name:      userPo.Name,
				Avatar:    userPo.Avatar,
				Mobile:    userPo.Mobile,
				Email:     userPo.Email,
				LoginName: userPo.LoginName,
			})
		}
		// 拼装用户参数
		userDeptIds, err := domain.GetUserDeptIds(orgID, userIds)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		deptPathNodeMap := map[int64][]resp.ContactPathNode{}
		for i, contactUser := range contactUserList {
			if deptIds, ok := userDeptIds[contactUser.ID]; ok && len(deptIds) > 0 {
				pathNodesList := make([][]resp.ContactPathNode, 0)
				for _, deptId := range deptIds {
					if pathNodes, exist := deptPathNodeMap[deptId]; exist {
						pathNodesList = append(pathNodesList, pathNodes)
					} else {
						pathNodes, err = getDepPathNodes(depMap, deptId, true)
						if err != nil {
							logger.Error(err)
							continue
						}
						deptPathNodeMap[deptId] = pathNodes
						pathNodesList = append(pathNodesList, pathNodes)
					}
				}
				contactUserList[i].PathNodesList = pathNodesList
			}
		}
		return &resp.ContactSearchResp{
			User:    contactUserList,
			HasMore: hasMore,
		}, nil
	}

	return &resp.ContactSearchResp{}, nil
}
