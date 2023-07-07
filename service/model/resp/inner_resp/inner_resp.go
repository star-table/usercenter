package inner_resp

type UserBaseInfo struct {
	Id            int64  `json:"id"`            // Id 用户ID
	Name          string `json:"name"`          // Name 姓名
	Avatar        string `json:"avatar"`        // Avatar 用户头像
	Email         string `json:"email"`         // Email 邮箱
	Mobile        string `json:"mobile"`        // Mobile 手机
	SourceChannel string `json:"sourceChannel"` // SourceChannel 渠道
	OpenId        string `json:"openId"`        // OpenId 第三方平台(飞书/钉钉/企微) OpenId
}
