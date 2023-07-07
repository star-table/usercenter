package open_req

// CompareUser
type CompareUser struct {
	UserId     int64 `json:"userId"`     // UserId
	SuperiorId int64 `json:"superiorId"` // SuperiorId 上级ID
}

// SameSubordinateUsers
type SameSubordinateUsers struct {
	Id      int64   `json:"id"`      // Id
	DeptIds []int64 `json:"deptIds"` // DeptIds
}
