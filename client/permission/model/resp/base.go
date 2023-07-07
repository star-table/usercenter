package resp

type Result struct {
	Message string `json:"message"`
	Code int `json:"code"`
	Data interface{} `json:"data"`
}