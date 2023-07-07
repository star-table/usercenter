package req

import "github.com/star-table/usercenter/service/model/bo"

/**
职级相关的响应数据模型
*/

// CreatePositionReq 创建职级请求数据模型
type CreatePositionReq struct {
	// Name 名称
	Name string `json:"name"`
	// PositionLevel 等级
	PositionLevel int `json:"positionLevel"`
	// Remark 说明
	Remark string `json:"remark"`
	// Status 状态
	Status int `json:"status"`
}

// ModifyPositionInfoReq 修改职级请求数据模型
type ModifyPositionInfoReq struct {
	// PositionId 组织内职级ID
	PositionId int64 `json:"positionId"`
	// Name 名称
	Name string `json:"name"`
	// PositionLevel 等级
	PositionLevel int `json:"positionLevel"`
	// Remark 说明
	Remark string `json:"remark"`
}

// UpdatePositionStatusReq 修改职级状态
type UpdatePositionStatusReq struct {
	// PositionId 组织内职级ID
	PositionId int64 `json:"positionId"`
	// Status 状态 1启用 2停用
	Status int `json:"status"`
}

// SearchPositionListReq 获取职级列表
type SearchPositionListReq struct {
	// Status 状态 1启用 2停用
	Status int `json:"status"`
}

// SearchPositionPageListReq 获取职级分页列表
type SearchPositionPageListReq struct {
	// Status 状态 1启用 2停用
	Status int `json:"status"`
	// 分页参数
	bo.PageBo
}
