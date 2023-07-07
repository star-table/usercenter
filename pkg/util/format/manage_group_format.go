package format

import (
	"regexp"
	"strconv"
)

var (
	manageGroupNameReg = regexp.MustCompile(ManageGroupNamePattern)
)

// VerifyManageGroupNameFormat 验证管理员组名称规范
func VerifyManageGroupNameFormat(input string) bool {
	formInput := ChineseReg.ReplaceAllString(input, "a")
	return manageGroupNameReg.MatchString(formInput)
}

var (
	AppOrPkgIdReg = regexp.MustCompile(`^(?P<type>[a|p])(?P<id>[0-9]+)$`)
)

// VerifyAndGetAppOrPkgIdFormat 验证PkgId规范
func VerifyAndGetAppOrPkgIdFormat(input string) (string, int64, bool) {
	match := AppOrPkgIdReg.FindStringSubmatch(input)
	if len(match) != 3 {
		return "", int64(0), false
	}
	id, err := strconv.ParseInt(match[2], 10, 64)
	if err != nil {
		return "", int64(0), false
	}
	return match[1], id, true
}
