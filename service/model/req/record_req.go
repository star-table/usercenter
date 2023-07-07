package req

import (
	"time"

	"github.com/star-table/usercenter/service/model/bo"
)

// LoginRecordListReq
type LoginRecordListReq struct {
	IP         string    `json:"ip"`        // IP
	AccountNo  string    `json:"accountNo"` // AccountNo
	StartDate  time.Time `json:"startDate"` // StartDate
	EndDate    time.Time `json:"endDate"`   // EndDate
	*bo.PageBo           // 分页
}

// ExportLoginRecordListReq
type ExportLoginRecordListReq struct {
	IP        string    `json:"ip"`        // IP
	AccountNo string    `json:"accountNo"` // AccountNo
	StartDate time.Time `json:"startDate"` // StartDate
	EndDate   time.Time `json:"endDate"`   // EndDate
}
