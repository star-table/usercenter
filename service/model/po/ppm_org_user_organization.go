package po

import "time"

type PpmOrgUserOrganization struct {
	Id               int64     `db:"id,omitempty" json:"id"`
	OrgId            int64     `db:"org_id,omitempty" json:"orgId"`
	UserId           int64     `db:"user_id,omitempty" json:"userId"`
	InviteId         int64     `db:"invite_id,omitempty" json:"inviteId"`
	CheckStatus      int       `db:"check_status,omitempty" json:"checkStatus"`
	UseStatus        int       `db:"use_status,omitempty" json:"useStatus"`
	Status           int       `db:"status,omitempty" json:"status"`
	Type             int       `db:"type,omitempty" json:"type"`
	EmpNo            string    `db:"emp_no,omitempty" json:"empNo"`       // 工号
	WeiboIds         string    `db:"weibo_ids,omitempty" json:"weiboIds"` // 微博ID array string ,号分割
	StatusChangerId  int64     `db:"status_changer_id,omitempty" json:"statusChangerId"`
	StatusChangeTime time.Time `db:"status_change_time,omitempty" json:"statusChangeTime"`
	AuditorId        int64     `db:"auditor_id,omitempty" json:"auditorId"`
	AuditTime        time.Time `db:"audit_time,omitempty" json:"auditTime"`
	Creator          int64     `db:"creator,omitempty" json:"creator"`
	CreateTime       time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator          int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime       time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version          int       `db:"version,omitempty" json:"version"`
	IsDelete         int       `db:"is_delete,omitempty" json:"isDelete"`
}

func (*PpmOrgUserOrganization) TableName() string {
	return "ppm_org_user_organization"
}

type PpmOrgUserOrganizationCount struct {
	Total uint64 `db:"total,omitempty" json:"total"`
}

type PpmUserId struct {
	UserId int64 `db:"user_id,omitempty" json:"userId"` // UserId
}
