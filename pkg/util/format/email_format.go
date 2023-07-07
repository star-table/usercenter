package format

import "regexp"

var (
	EmailReg = regexp.MustCompile(EmailPattern)
)

// VerifyEmailFormat 邮箱认证
func VerifyEmailFormat(email string) bool {
	if SpaceReg.MatchString(email) {
		return false
	}
	if !EmailReg.MatchString(email) {
		return false
	}
	return true
}
