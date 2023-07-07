package resp

import "time"

type CreateOrgRespVoData struct {
	OrgId int64 `json:"orgId"`
}

type UserOrganizationListResp struct {
	// 用户组织列表
	List []UserOrganization `json:"list"`
}

// UserOrganization 用户组织列表响应结构体
type UserOrganization struct {
	// 组织id
	ID int64 `json:"id"`
	// 组织名称
	Name string `json:"name"`
	// 组织code
	Code string `json:"code"`
	// 组织网站
	WebSite string `json:"webSite"`
	// 所属行业
	IndustryID int64 `json:"industryId"`
	// 组织规模
	Scale string `json:"scale"`
	// 来源平台
	SourcePlatform string `json:"sourcePlatform"`
	// 来源渠道
	SourceChannel string `json:"sourceChannel"`
	// 所在国家
	CountryID int64 `json:"countryId"`
	// 所在省份
	ProvinceID int64 `json:"provinceId"`
	// 所在城市
	CityID int64 `json:"cityId"`
	// 组织地址
	Address string `json:"address"`
	// 组织logo地址
	LogoURL string `json:"logoUrl"`
	// 组织标识
	ResorceID int64 `json:"resorceId"`
	// 组织所有人,创建时默认为创建人
	Owner int64 `json:"owner"`
	// 企业是否认证
	IsAuthenticated int `json:"IsAuthenticated"`
	// 描述
	Remark string `json:"remark"`
	// 是否展示
	IsShow int `json:"isShow"`
	// 是否删除,1是,2否
	IsDelete int `json:"isDelete"`
	// 对于该用户组织是否可用（1是2否）
	OrgIsEnabled int `json:"OrgIsEnabled"`
	// 组织可用功能
	Functions []string `json:"functions"`
}

type OrgConfig struct {
	Id                         int64     `db:"id,omitempty" json:"id"`
	OrgId                      int64     `db:"org_id,omitempty" json:"orgId"`
	TimeZone                   string    `db:"time_zone,omitempty" json:"timeZone"`
	TimeDifference             string    `db:"time_difference,omitempty" json:"timeDifference"`
	DbId                       int64     `db:"db_id,omitempty" json:"dbId"`
	DcId                       int64     `db:"dc_id,omitempty" json:"dcId"`
	DsId                       int64     `db:"ds_id,omitempty" json:"dsId"`
	PayLevel                   int       `db:"pay_level,omitempty" json:"payLevel"`
	PayStartTime               time.Time `db:"pay_start_time,omitempty" json:"payStartTime"`
	PayEndTime                 time.Time `db:"pay_end_time,omitempty" json:"payEndTime"`
	WebSite                    string    `db:"web_site,omitempty" json:"webSite"`
	Language                   string    `db:"language,omitempty" json:"language"`
	RemindSendTime             string    `db:"remind_send_time,omitempty" json:"remindSendTime"`
	ProjectDailyReportSendTime string    `db:"project_daily_report_send_time,omitempty" json:"projectDailyReportSendTime"`
	DatetimeFormat             string    `db:"datetime_format,omitempty" json:"datetimeFormat"`
	PasswordLength             int       `db:"password_length,omitempty" json:"passwordLength"`
	PasswordRule               int       `db:"password_rule,omitempty" json:"passwordRule"`
	MaxLoginFailCount          int       `db:"max_login_fail_count,omitempty" json:"maxLoginFailCount"`
	Status                     int       `db:"status,omitempty" json:"status"`
	Creator                    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime                 time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator                    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime                 time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version                    int       `db:"version,omitempty" json:"version"`
	IsDelete                   int       `db:"is_delete,omitempty" json:"isDelete"`
}
