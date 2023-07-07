package req

type AppListReq struct {
	OrgId   int64 `json:"orgId,omitempty"`
	GroupId int64 `json:"groupId,omitempty"`
	PkgId   int64 `json:"pkgId,omitempty"`
	Type    int64 `json:"type,omitempty"`
}
