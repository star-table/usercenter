package resp

type GetOssPostPolicyResp struct {
	// policy
	Policy string `json:"policy"`
	// 签名
	Signature string `json:"signature"`
	// 文件上传目录
	Dir string `json:"dir"`
	// 有效期
	Expire string `json:"expire"`
	// access Id
	AccessID string `json:"accessId"`
	// Host
	Host string `json:"host"`
	// Region
	Region string `json:"region"`
	// bucket名称
	Bucket string `json:"bucket"`
	// 文件名
	FileName string `json:"fileName"`
	// 文件最大限制
	MaxFileSize int64 `json:"maxFileSize"`
	// callback回调，为空说明不需要回调
	Callback string `json:"callback"`
}
