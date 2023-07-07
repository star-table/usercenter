package consts

import "github.com/star-table/usercenter/pkg/util/slice"

const (
	// 部门负责人
	DepartmentIsLeader = 1
	// 非部门负责人
	DepartmentNotLeader = 2
)

var (
	leaderValues = []int{DepartmentIsLeader, DepartmentNotLeader}
)

func CheckLeaderValue(value int) bool {
	ok, _ := slice.Contain(leaderValues, value)
	return ok
}
