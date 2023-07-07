package bo

import (
	"time"

	"github.com/star-table/usercenter/pkg/store/mysql"
)

type OrganizationBo struct {
	Id              int64     `db:"id,omitempty" json:"id"`
	Name            string    `db:"name,omitempty" json:"name"`
	Code            string    `db:"code,omitempty" json:"code"`
	WebSite         string    `db:"web_site,omitempty" json:"webSite"`
	IndustryId      int64     `db:"industry_id,omitempty" json:"industryId"`
	Scale           string    `db:"scale,omitempty" json:"scale"`
	SourceChannel   string    `db:"source_channel,omitempty" json:"sourceChannel"`
	CountryId       int64     `db:"country_id,omitempty" json:"countryId"`
	ProvinceId      int64     `db:"province_id,omitempty" json:"provinceId"`
	CityId          int64     `db:"city_id,omitempty" json:"cityId"`
	Address         string    `db:"address,omitempty" json:"address"`
	LogoUrl         string    `db:"logo_url,omitempty" json:"logoUrl"`
	ResourceId      int64     `db:"resource_id,omitempty" json:"resourceId"`
	Owner           int64     `db:"owner,omitempty" json:"owner"`
	IsAuthenticated int       `db:"is_authenticated,omitempty" json:"isAuthenticated"`
	Remark          string    `db:"remark,omitempty" json:"remark"`
	InitStatus      int       `db:"init_status,omitempty" json:"initStatus"`
	InitVersion     int       `db:"init_version,omitempty" json:"initVersion"`
	Status          int       `db:"status,omitempty" json:"status"`
	IsShow          int       `db:"is_show,omitempty" json:"isShow"`
	ApiKey          string    `db:"api_key" json:"apiKey"`
	Creator         int64     `db:"creator,omitempty" json:"creator"`
	CreateTime      time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator         int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime      time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version         int       `db:"version,omitempty" json:"version"`
	IsDelete        int       `db:"is_delete,omitempty" json:"isDelete"`
}

type ScheduleOrganizationListBo struct {
	OrgId                      int64  `json:"orgId"`
	ProjectDailyReportSendTime string `json:"projectDailyReportSendTime"`
}

type CreateOrgBo struct {
	OrgName    string  `json:"orgName"`
	IndustryId *int64  `json:"industryId"`
	Scale      *string `json:"scale"`
}

type UpdateOrganizationBo struct {
	Bo                     OrganizationBo
	OrganizationUpdateCond mysql.Upd
}

type OrgIdInfo struct {
	OrgId    int64  `json:"orgId"`
	OutOrgId string `json:"outOrgId"`
}
