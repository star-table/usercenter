package inner_resp

type GetManagerResp struct {
	Data []GetManagerData `json:"data"`
}

type GetManagerData struct {
	MemberType string  `json:"memberType"`
	MemberId   int64   `json:"memberId"`
	LangCode   string  `json:"langCode"`
	AppIds     []int64 `json:"appIds"`
	IsSysAdmin bool    `json:"isSysAdmin"`
	IsSubAdmin bool    `json:"isSubAdmin"`
}

type ManagerGroupInitResp struct {
	SysGroupID int64 `json:"sysGroupId,string"`
}
