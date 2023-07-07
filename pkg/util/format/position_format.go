package format

import "regexp"

var (
	positionNameReg = regexp.MustCompile(PositionNamePattern)
)

/**
VerifyPositionNameFormat 验证职级名称规范
@author WangShiChang
@version v1.0
@date 2020-10-21
*/
func VerifyPositionNameFormat(input string) bool {
	// 过滤空格和中文
	if SpaceReg.MatchString(input) {
		return false
	}
	formInput := ChineseReg.ReplaceAllString(input, "a")
	return positionNameReg.MatchString(formInput)
}
