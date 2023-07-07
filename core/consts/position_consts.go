package consts

/**
管理组
*/

const (
	PositionManagerId = int64(1)
	PositionMemberId  = int64(2)
)
const (
	PositionManagerName = "主管"
	PositionMemberName  = "成员"
)
const (
	PositionManagerLevel = 10
	PositionMemberLevel  = 20
)

// DefaultPositions 默认职级名称和级别
var DefaultPositions = map[int64]map[string]interface{}{
	PositionManagerId: {
		TcName:          PositionManagerName,
		TcPositionLevel: PositionManagerLevel,
	},
	PositionMemberId: {
		TcName:          PositionMemberName,
		TcPositionLevel: PositionMemberLevel,
	},
}
