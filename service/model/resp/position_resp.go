package resp

import (
	"time"
)

/**
职级信息响应数据模型
@author WangShiChang
@version v1.0
@date 2020-10-21
*/

// PositionInfoResp 职级信息响应数据模型
type PositionInfoResp struct {
	Id            int64     `json:"id"`
	PositionId    int64     `json:"positionId"`
	OrgId         int64     `json:"orgId"`
	Name          string    `json:"name"`
	PositionLevel int       `json:"positionLevel"`
	Remark        string    `json:"remark"`
	Status        int       `json:"status"`
	Creator       int64     `json:"creator"`
	CreateTime    time.Time `json:"createTime"`
	Updator       int64     `json:"updator"`
	UpdateTime    time.Time `json:"updateTime"`
	IsDelete      int       `json:"isDelete"`
	IsDefault     bool      `json:"isDefault"` //是否是默认的
}

// PositionPageListResp 职级分页列表响应数据模型
type PositionPageListResp struct {
	Total int64              `json:"total"`
	List  []PositionInfoResp `json:"list"`
}
