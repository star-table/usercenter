package bo

import "time"

type ManageGroupBo struct {
	Id            int64     `db:"id,omitempty" json:"id"`
	OrgId         int64     `db:"org_id,omitempty" json:"orgId"`
	Name          string    `db:"name,omitempty" json:"name"`
	LangCode      string    `db:"lang_code,omitempty" json:"langCode"`
	UserIds       string    `db:"user_ids,omitempty" json:"userIds"`    // 成员
	AppPackageIds string    `db:"app_package_ids" json:"appPackageIds"` // 应用包权限
	AppIds        string    `db:"app_ids" json:"appIds"`                // 应用权限
	UsageIds      string    `db:"usage_ids,omitempty" json:"usageIds"`  // 功能权限
	DeptIds       string    `db:"dept_ids,omitempty" json:"deptIds"`    // 部门权限
	RoleIds       string    `db:"role_ids,omitempty" json:"roleIds"`    // 角色权限
	Creator       int64     `db:"creator,omitempty" json:"creator"`
	CreateTime    time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator       int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime    time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version       int       `db:"version,omitempty" json:"version"`
	IsDelete      int       `db:"is_delete,omitempty" json:"isDelete"`
}

type ManageGroupInfoBo struct {
	Id            int64  `db:"id,omitempty" json:"id"`
	OrgId         int64  `db:"org_id,omitempty" json:"orgId"`
	Name          string `db:"name,omitempty" json:"name"`
	LangCode      string `db:"lang_code,omitempty" json:"langCode"`
	UserIds       string `db:"user_ids,omitempty" json:"userIds"`    // 成员
	AppPackageIds string `db:"app_package_ids" json:"appPackageIds"` // 应用包权限
	AppIds        string `db:"app_ids" json:"appIds"`                // 应用权限
	UsageIds      string `db:"usage_ids,omitempty" json:"usageIds"`  // 功能权限
	DeptIds       string `db:"dept_ids,omitempty" json:"deptIds"`    // 部门权限
	RoleIds       string `db:"role_ids,omitempty" json:"roleIds"`    // 角色权限
}
