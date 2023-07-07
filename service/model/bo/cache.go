package bo

//缓存用户登录信息
type CacheUserInfoBo struct {
	OutUserId     string `json:"outUserId"`
	SourceChannel string `json:"sourceChannel"`
	UserId        int64  `json:"userId"`
	CorpId        string `json:"corpId"`
	OrgId         int64  `json:"orgId"`
}

type CacheFsUserTokenBo struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
	FetchTime    int64  `json:"fetchTime"`
}

type BaseUserInfoBo struct {
	UserId        int64  `json:"userId"`
	OutUserId     string `json:"outUserId"` //有可能为空
	OrgId         int64  `json:"orgId"`
	OutOrgId      string `json:"outOrgId"`  //有可能为空
	LoginName     string `json:"loginName"` // 用户名
	Name          string `json:"name"`      // 昵称
	NamePy        string `json:"namePy"`
	Email         string `json:"email"`
	MobileRegion  string `json:"mobileRegion"` // MobileRegion
	Mobile        string `json:"mobile"`
	Avatar        string `json:"avatar"`
	HasOutInfo    bool   `json:"hasOutInfo"`
	HasOrgOutInfo bool   `json:"hasOrgOutInfo"`
	OutOrgUserId  string `json:"outOrgUserId"` //可能为空（钉钉发送消息会用到）

	OrgUserIsDelete    int `json:"orgUserIsDelete"`    //是否被组织移除
	OrgUserStatus      int `json:"orgUserStatus"`      //用户组织状态
	OrgUserCheckStatus int `json:"orgUserCheckStatus"` //用户组织审核状态

	EmpNo    string   `json:"empNo"`    // 工号
	WeiboIds []string `json:"weiboIds"` // 微博ID
}

//用户基本信息扩展
type BaseUserInfoExtBo struct {
	BaseUserInfoBo

	//部门id
	DepartmentId int64 `json:"departmentId"`
}

type BaseUserOutInfoBo struct {
	UserId       int64  `json:"userId"`
	OutUserId    string `json:"outUserId"` //有可能为空
	OutOrgId     string `json:"outOrgId"`  //有可能为空
	OrgId        int64  `json:"orgId"`
	OutOrgUserId string `json:"ourtOrgUserId"` //有可能为空
}

type BaseOrgInfoBo struct {
	OrgId         int64  `json:"orgId"`
	OrgName       string `json:"orgName"`
	OrgOwnerId    int64  `json:"orgOwnerId"`
	Creator       int64  `json:"creator"`
	OutOrgId      string `json:"outOrgId"`
	PayLevel      int    `json:"payLevel"`
	SourceChannel string `json:"sourceChannel"`
}

type BaseOrgOutInfoBo struct {
	OrgId         int64  `json:"orgId"`
	OutOrgId      string `json:"outOrgId"`
	SourceChannel string `json:"sourceChannel"`
}

type CacheProcessStatusBo struct {
	StatusId    int64  `json:"statusId"`
	StatusType  int    `json:"statusType"`
	Category    int    `json:"category"`
	IsInit      bool   `json:"isInit"`
	BgStyle     string `json:"bgStyle"`
	FontStyle   string `json:"fontStyle"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	LangCode    string `json:"langCode"`
	Sort        int    `json:"sort"`
	ProcessId   int64  `json:"processId"`
}

type CacheProjectCalendarInfoBo struct {
	IsSyncOutCalendar int    `json:"isSyncOutCalendar"`
	CalendarId        string `json:"calendarId"`
}
