package inner_req

type GetBaseInfoByIdsReq struct {
	OrgId int64    `json:"orgId"`
	Ids   []string `json:"ids"`
}
