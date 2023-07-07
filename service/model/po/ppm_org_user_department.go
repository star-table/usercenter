package po

import "time"

type PpmOrgUserDepartment struct {
	Id            int64     `db:"id,omitempty" json:"id"`
	OrgId         int64     `db:"org_id,omitempty" json:"orgId"`
	UserId        int64     `db:"user_id,omitempty" json:"userId"`
	DepartmentId  int64     `db:"department_id,omitempty" json:"departmentId"`
	IsLeader      int       `db:"is_leader,omitempty" json:"isLeader"`            // 是否是部门负责人 1是 2否
	OrgPositionId int64     `db:"org_position_id,omitempty" json:"orgPositionId"` //注意此处是org内的职级ID
	Creator       int64     `db:"creator,omitempty" json:"creator"`
	CreateTime    time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator       int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime    time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version       int       `db:"version,omitempty" json:"version"`
	IsDelete      int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmOrgUserDepartment) TableName() string {
	return "ppm_org_user_department"
}
