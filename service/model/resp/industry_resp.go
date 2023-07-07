package resp

type IndustryListResp struct {
	List []*IndustryResp `json:"list"`
}

type IndustryResp struct {
	// 主键
	ID int64 `json:"id"`
	// 名字
	Name string `json:"name"`
	// 中文名
	Cname string `json:"cname"`
}
