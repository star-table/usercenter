package bo

import "time"

// OrgMemberBaseInfoBo 基础组织成员信息
type OrgMemberBaseInfoBo struct {
	OrgId            int64     `db:"org_id,omitempty" json:"orgId"`
	OrgName          string    `db:"org_name,omitempty" json:"orgName"` // OrgName
	UserId           int64     `db:"user_id,omitempty" json:"userId"`
	LoginName        string    `db:"login_name,omitempty" json:"loginName"`
	Name             string    `db:"name,omitempty" json:"name"`
	NamePinyin       string    `db:"name_pinyin,omitempty" json:"namePinyin"`
	Email            string    `db:"email,omitempty" json:"email"`
	MobileRegion     string    `db:"mobile_region,omitempty" json:"mobileRegion"`
	Mobile           string    `db:"mobile,omitempty" json:"mobile"`
	Sex              string    `db:"sex,omitempty" json:"sex"`           // Sex
	Birthday         time.Time `db:"birthday,omitempty" json:"birthday"` // Birthday
	Language         string    `db:"language,omitempty" json:"language"` // Language
	Avatar           string    `db:"avatar,omitempty" json:"avatar"`
	EmpNo            string    `db:"emp_no,omitempty" json:"empNo"`       // 工号
	WeiboIds         string    `db:"weibo_ids,omitempty" json:"weiboIds"` // 微博ID array string ,号分割
	CheckStatus      int       `db:"check_status,omitempty" json:"checkStatus"`
	Status           int       `db:"status,omitempty" json:"status"`
	Type             int       `db:"type,omitempty" json:"type"`
	StatusChangeTime time.Time `db:"status_change_time,omitempty" json:"statusChangeTime"` // status修改时间
	UseStatus        int       `db:"use_status,omitempty" json:"useStatus"`                // 使用过该组织
	AuditorId        int64     `db:"auditor_id,omitempty" json:"auditorId"`                // 审核人
	AuditTime        time.Time `db:"audit_time,omitempty" json:"auditTime"`                // 审核时间
	OrgCreator       int64     `db:"org_creator,omitempty" json:"orgCreator"`              // 组织创建者
	OrgOwner         int64     `db:"org_owner,omitempty" json:"orgOwner"`                  // 组织拥有者
	Creator          int64     `db:"creator,omitempty" json:"creator"`
	CreateTime       time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator          int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime       time.Time `db:"update_time,omitempty" json:"updateTime"`
	IsDelete         int       `db:"is_delete,omitempty" json:"isDelete"` // IsDelete
}

type OrgUserSimpleInfoBo struct {
	Id     int64  `db:"id,omitempty" json:"id"`
	Name   string `db:"name,omitempty" json:"name"`
	Avatar string `db:"avatar,omitempty" json:"avatar"`
	Status int    `db:"status,omitempty" json:"status"`
}

type UserBaseInfoBo struct {
	Id            int64   `db:"id,omitempty" json:"id"`                        // Id 用户ID
	Name          string  `db:"name,omitempty" json:"name"`                    // Name 姓名
	Avatar        string  `db:"avatar,omitempty" json:"avatar"`                // Avatar 用户头像
	Email         string  `db:"email,omitempty" json:"email"`                  // Email 邮箱
	Mobile1       string  `db:"mobile1,omitempty" json:"mobile1"`              // Mobile 手机
	Mobile2       *string `db:"mobile2,omitempty" json:"mobile2"`              // Mobile 手机
	SourceChannel string  `db:"source_channel,omitempty" json:"sourceChannel"` // SourceChannel 渠道
	OpenId        string  `db:"open_id,omitempty" json:"openId"`               // OpenId 第三方平台(飞书/钉钉/企微) OpenId
}
