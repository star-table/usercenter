package resp

type AppPackageListReq struct {
	Code      int              `json:"code"`
	Data      []AppPackageData `json:"data"`
	Message   string           `json:"message"`
	Success   bool             `json:"success"`
	Timestamp int64            `json:"timestamp"`
}

type AppPackageData struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type AppListReq struct {
	Code      int       `json:"code"`
	Data      []AppData `json:"data"`
	Message   string    `json:"message"`
	Success   bool      `json:"success"`
	Timestamp int64     `json:"timestamp"`
}

type AppData struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
