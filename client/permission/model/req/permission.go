package req

type InitDefaultManageGroupReq struct {
	OrgID       int64               `json:"orgId"`
	AuthOptions []OptAuthOptionInfo `json:"authOptions"`
}

type OptAuthOptionInfo struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Group    string `json:"group"`
	Required bool   `json:"required"`
	IsMenu   bool   `json:"isMenu"`
	Status   int    `json:"status"`
}

type FieldAuthOptionInfo struct {
	Code     int    `json:"code"`
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

type InitAppPermissionReq struct {
	OrgID                      int64                 `json:"orgId"`
	AppPackageID               int64                 `json:"appPackageId"`
	AppID                      int64                 `json:"appId"`
	AppType                    int                   `json:"appType"`
	OptAuthOptions             []OptAuthOptionInfo   `json:"optAuthOptions"`
	FieldAuthOptions           []FieldAuthOptionInfo `json:"fieldAuthOptions"`
	IsExt                      bool                  `json:"isExt"`
	ComponentType              string                `json:"componentType"`
	Creatable                  bool                  `json:"creatable"`
	UserID                     int64                 `json:"userId"`
	Config                     string                `json:"config"`
	DefaultPermissionGroupType int                   `json:"defaultPermissionGroupType"`
}
