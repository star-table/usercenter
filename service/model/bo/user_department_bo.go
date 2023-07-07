package bo

import "time"

type UserDepartmentBo struct {
	Id           int64     `db:"id,omitempty" json:"id"`
	OrgId        int64     `db:"org_id,omitempty" json:"orgId"`
	UserId       int64     `db:"user_id,omitempty" json:"userId"`
	DepartmentId int64     `db:"department_id,omitempty" json:"departmentId"`
	IsLeader     int       `db:"is_leader,omitempty" json:"isLeader"`
	Creator      int64     `db:"creator,omitempty" json:"creator"`
	CreateTime   time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator      int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime   time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version      int       `db:"version,omitempty" json:"version"`
	IsDelete     int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*UserDepartmentBo) TableName() string {
	return "ppm_org_user_department"
}

type DepartUserCount struct {
	DepartmentID int64 `json:"departmentId" db:"department_id"`
	Count uint64 `json:"count" db:"count"`
}

type DepartUserIds struct {
	DepartmentID int64 `json:"departmentId" db:"department_id"`
	UserIds uint64 `json:"userIds" db:"user_ids"`
}

type UserDepts struct {
	UserID int64 `json:"userId"`
	DeptIds []int64 `json:"deptIds"`
}