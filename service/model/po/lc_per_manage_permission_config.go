package po

import "time"

type LcPerManagePermissionConfig struct {
	Id             int64     `db:"id,omitempty" json:"id"`
	OrgId          int64     `db:"org_id,omitempty" json:"orgId"`
	OptAuthOptions string    `db:"opt_auth_options,omitempty" json:"optAuthOptions"`
	Type           int       `db:"type,omitempty" json:"type"`
	Creator        int64     `db:"creator,omitempty" json:"creator"`
	CreateTime     time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator        int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime     time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version        int       `db:"version,omitempty" json:"version"`
	DelFlag        int       `db:"del_flag,omitempty" json:"delFlag"`
}

func (*LcPerManagePermissionConfig) TableName() string {
	return "lc_per_manage_permission_config"
}

