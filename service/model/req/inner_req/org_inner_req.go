package inner_req

type AddOrgOutCollaborator struct {
	OrgID  int64 `json:"orgId"`
	UserID int64 `json:"userId"`
}

type CheckAndSetSuperAdminReq struct {
	OrgID  int64 `json:"orgId"`
	UserID int64 `json:"userId"`
}
