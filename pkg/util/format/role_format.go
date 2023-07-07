package format

import "regexp"

var (
	roleNameReg      = regexp.MustCompile(RoleNamePattern)
	roleGroupNameReg = regexp.MustCompile(RoleGroupNamePattern)
)

// VerifyRoleNameFormat 验证角色名称规范
func VerifyRoleNameFormat(input string) bool {
	formInput := ChineseReg.ReplaceAllString(input, "a")
	return roleNameReg.MatchString(formInput)
}

// VerifyRoleGroupNameFormat 验证角色组名称规范
func VerifyRoleGroupNameFormat(input string) bool {
	formInput := ChineseReg.ReplaceAllString(input, "a")
	return roleGroupNameReg.MatchString(formInput)
}
