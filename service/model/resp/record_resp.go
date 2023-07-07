package resp

import (
	"time"
)

// LoginRecordListResp
type LoginRecordListResp struct {
	List  []LoginRecordData `json:"list"`  // List
	Total int64             `json:"total"` // Total
}

// LoginRecordData
type LoginRecordData struct {
	ID         int64     `json:"id"`         // ID
	AccountNo  string    `json:"accountNo"`  // AccountNo
	UserId     int64     `json:"userId"`     // UserId
	IP         string    `json:"ip"`         // IP
	UserAgent  string    `json:"userAgent"`  // UserAgent
	Msg        string    `json:"msg"`        // Msg
	CreateTime time.Time `json:"createTime"` // CreateTime
}
