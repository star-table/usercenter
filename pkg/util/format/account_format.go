package format

import (
	"regexp"
)

var (
	// AccountPwdReg 账号密码正则验证器
	AccountPwdReg = regexp.MustCompile(PasswordPattern)
)

// VerifyAccountPwdFormat 验证账号密码
func VerifyAccountPwdFormat(input string) bool {
	// 过滤空格和中文
	if SpaceReg.MatchString(input) || ChineseReg.MatchString(input) {
		return false
	}

	pwdLen := len(input)
	if pwdLen < 8 || pwdLen > 16 {
		return false
	}

	if !UpperCaseLetterReg.MatchString(input) {
		if !LowerCaseLetterReg.MatchString(input) {
			return false
		}
		if !NumberReg.MatchString(input) {
			return false
		}
		if len(NumberReg.ReplaceAllString(LowerCaseLetterReg.ReplaceAllString(input, ""), "")) > 0 {
			return true
		}
		return false
	} else if !LowerCaseLetterReg.MatchString(input) {
		if !UpperCaseLetterReg.MatchString(input) {
			return false
		}
		if !NumberReg.MatchString(input) {
			return false
		}
		if len(NumberReg.ReplaceAllString(UpperCaseLetterReg.ReplaceAllString(input, ""), "")) > 0 {
			return true
		}
		return false
	} else if !NumberReg.MatchString(input) {
		if !UpperCaseLetterReg.MatchString(input) {
			return false
		}
		if !LowerCaseLetterReg.MatchString(input) {
			return false
		}
		if len(LowerCaseLetterReg.ReplaceAllString(UpperCaseLetterReg.ReplaceAllString(input, ""), "")) > 0 {
			return true
		}
		return false
	}

	return true
}

// VerifyNicknameFormat 昵称验证
func VerifyNicknameFormat(input string) bool {
	if SpaceReg.MatchString(input) {
		return false
	}
	reg := regexp.MustCompile(ChinesePattern)
	formInput := reg.ReplaceAllString(input, "a")
	return len(formInput) <= 20
}

var AccountReg = regexp.MustCompile(AccountPattern)

// VerifyAccountFormat 账号验证
func VerifyAccountFormat(input string) bool {
	// 过滤空格和中文
	if SpaceReg.MatchString(input) || ChineseReg.MatchString(input) {
		return false
	}
	return AccountReg.MatchString(input)
}

var (
	// 手机区号
	MobileRegionReg = regexp.MustCompile(MobileRegionPattern)
	MobileReg       = regexp.MustCompile(MobilePattern)
)

// VerifyMobileRegionFormat 手机区号验证
func VerifyMobileRegionFormat(input string) bool {
	return MobileRegionReg.MatchString(input)
}

// VerifyMobileFormat 手机号验证
func VerifyMobileFormat(input string) bool {
	return MobileReg.MatchString(input)
}
