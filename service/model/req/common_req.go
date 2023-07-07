package req

// IdReq Id 请求 数据模型
type IdReq struct {
	// ID
	Id int64 `json:"id"`
}

// IdsReq Ids 请求 数据模型
type IdsReq struct {
	// Ids ID列表
	Ids []int64 `json:"ids"`
}
