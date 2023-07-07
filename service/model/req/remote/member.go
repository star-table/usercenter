package remote

import "time"

type UserData struct {
	Id          int64     `json:"id"`          // Id 用户ID
	LoginName   string    `json:"loginName"`   // LoginName
	Name        string    `json:"name"`        // Name 姓名
	NamePy      string    `json:"namePy"`      // NamePy 姓名拼音
	Avatar      string    `json:"avatar"`      // Avatar 用户头像
	Email       string    `json:"email"`       // Email 邮箱
	PhoneRegion string    `json:"phoneRegion"` // PhoneRegion 手机区号
	PhoneNumber string    `json:"phoneNumber"` // PhoneNumber 手机
	Status      int       `json:"status"`      // Status 状态1启用2禁用3离职
	Creator     int64     `json:"creator"`     // Creator 组织成员 创建者
	CreateTime  time.Time `json:"createTime"`  // CreateTime 创建时间
	Updator     int64     `json:"updator"`     // Updator 组织成员 最后修改者ID
	UpdateTime  time.Time `json:"updateTime"`  // UpdateTime 创建时间
}

type UserList struct {
	UserList []UserData `json:"userList"` // UserList
}
