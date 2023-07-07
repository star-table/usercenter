package resp

type AllocateDbResp struct {
	Code      int       `json:"code"`
	Data      AllocateData `json:"data"`
	Message   string    `json:"message"`
}

type AllocateData struct {
	DbId int64 `json:"dbId"`
	DsId int64 `json:"dsId"`
	DcId int64 `json:"dcId"`
}

