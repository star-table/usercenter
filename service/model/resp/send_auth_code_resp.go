package resp

type AuthSmsCodeResp struct {
	// 校验成功后的 token，通过此 token 可以做一些操作。有效期 5 min
	Token    string `json:"token"`
	// 验证方式。接口入参中传过来的。
	AuthType int    `json:"authType"`
}
