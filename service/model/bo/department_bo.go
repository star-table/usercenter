package bo

import "time"

// OrgDeptBo
type OrgDeptBo struct {
	Id                       int64     `db:"id,omitempty" json:"id"`
	OrgId                    int64     `db:"org_id,omitempty" json:"orgId"`
	Name                     string    `db:"name,omitempty" json:"name"`
	Code                     string    `db:"code,omitempty" json:"code"`
	ParentId                 int64     `db:"parent_id,omitempty" json:"parentId"`
	Path                     string    `db:"path,omitempty" json:"path"`
	Sort                     int       `db:"sort,omitempty" json:"sort"`
	IsHide                   int       `db:"is_hide,omitempty" json:"isHide"`
	SourcePlatform           string    `db:"source_platform,omitempty" json:"sourcePlatform"`
	SourceChannel            string    `db:"source_channel,omitempty" json:"sourceChannel"`
	OutOrgDepartmentId       string    `db:"out_org_department_id,omitempty" json:"outOrgDepartmentId"`
	OutOrgDepartmentCode     string    `db:"out_org_department_code,omitempty" json:"outOrgDepartmentCode"`
	OutOrgDepartmentParentId string    `db:"out_org_department_parent_id,omitempty" json:"outOrgDepartmentParentId"`
	Status                   int       `db:"status,omitempty" json:"status"`
	Creator                  int64     `db:"creator,omitempty" json:"creator"`
	CreateTime               time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator                  int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime               time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version                  int       `db:"version,omitempty" json:"version"`
	IsDelete                 int       `db:"is_delete,omitempty" json:"isDelete"`
}

// UserDeptBindBo 对应下方的指针版 `UserDeptBindBoValPtr`，这两个结构体字段名请保持一致。
type UserDeptBindBo struct {
	OrgId                    int64  `db:"org_id,omitempty" json:"orgId"`
	UserId                   int64  `db:"user_id,omitempty" json:"userId"`
	DepartmentId             int64  `db:"department_id,omitempty" json:"departmentId"`
	DepartmentName           string `db:"department_name,omitempty" json:"departmentName"`
	OutOrgDepartmentId       string `db:"out_org_department_id,omitempty" json:"outOrgDepartmentId"`
	OutOrgDepartmentCode     string `db:"out_org_department_code,omitempty" json:"outOrgDepartmentCode"`
	OutOrgDepartmentParentId string `db:"out_org_department_parent_id,omitempty" json:"outOrgDepartmentParentId"`
	IsLeader                 int    `db:"is_leader,omitempty" json:"isLeader"`            // 1是 2否
	OrgPositionId            int64  `db:"org_position_id,omitempty" json:"orgPositionId"` //注意此处是org内的职级ID
	PositionName             string `db:"position_name,omitempty" json:"positionName"`
	PositionLevel            int    `db:"position_level,omitempty" json:"positionLevel"`
}

type UserDeptBindBoValPtr struct {
	OrgId                    int64   `db:"org_id,omitempty" json:"orgId"`
	UserId                   *int64  `db:"user_id,omitempty" json:"userId"`
	DepartmentId             *int64  `db:"department_id,omitempty" json:"departmentId"`
	DepartmentName           *string `db:"department_name,omitempty" json:"departmentName"`
	OutOrgDepartmentId       *string `db:"out_org_department_id,omitempty" json:"outOrgDepartmentId"`
	OutOrgDepartmentCode     *string `db:"out_org_department_code,omitempty" json:"outOrgDepartmentCode"`
	OutOrgDepartmentParentId *string `db:"out_org_department_parent_id,omitempty" json:"outOrgDepartmentParentId"`
	IsLeader                 *int    `db:"is_leader,omitempty" json:"isLeader"`            // 1是 2否
	OrgPositionId            *int64  `db:"org_position_id,omitempty" json:"orgPositionId"` //注意此处是org内的职级ID
	PositionName             *string `db:"position_name,omitempty" json:"positionName"`
	PositionLevel            *int    `db:"position_level,omitempty" json:"positionLevel"`
}

type GetDeptUserIdsParams struct {
	Query   *string  `json:"query"`
	UserIds *[]int64 `json:"userIds"`
}

// 部门树节点
type DepartmentTreeNode struct {
	ID       int64                 `json:"id"`
	Name     string                `json:"name"`
	ParentID int64                 `json:"parentId"`
	Childs   []*DepartmentTreeNode `json:"childs"`
	Parent   *DepartmentTreeNode   `json:"parent"`
}

// 遍历节点，如果call返回false，则跳过当前结点的递归
func (d *DepartmentTreeNode) Foreach(call func(d *DepartmentTreeNode) bool) {
	for _, child := range d.Childs {
		if !call(child) {
			continue
		}
		child.Foreach(call)
	}
}

func (d *DepartmentTreeNode) Walk(f func(d *DepartmentTreeNode)) {
	f(d)
	for _, child := range d.Childs {
		child.Walk(f)
	}
}
