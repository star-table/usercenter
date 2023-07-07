package format

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestVerifyPwdFormat(t *testing.T) {
	suc := VerifyAccountPwdFormat("zZ12345，")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyAccountPwdFormat("Zz12345,")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyAccountPwdFormat("123dsf~~~")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyAccountPwdFormat("zZ123456.")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyAccountPwdFormat("123123qwe.")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyAccountPwdFormat("a123a123")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyAccountPwdFormat("a123a123'")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyAccountPwdFormat("1a")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyAccountPwdFormat("aAAAAA1")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyAccountPwdFormat("a")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyAccountPwdFormat("1")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyAccountPwdFormat("A1a&")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyAccountPwdFormat("a1a.")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyAccountPwdFormat("a%!#12")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyAccountPwdFormat("*")
	t.Log(suc)
	assert.Equal(t, suc, false)
}

func TestVerifyUserNameFormat(t *testing.T) {
	suc := VerifyNicknameFormat("好好好好好好好好好好")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyNicknameFormat("11111111111111111111")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyNicknameFormat("llllllllllllllllllll")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyNicknameFormat("好好好好好好好好好1K")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyNicknameFormat("好好好好好好好好好1Kl")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyNicknameFormat("*")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyNicknameFormat("hao好*")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyNicknameFormat("hao1265")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyNicknameFormat("hao1265             ")
	t.Log(suc)
	assert.Equal(t, suc, true)

	suc = VerifyNicknameFormat("hao1265              ")
	t.Log(suc)
	assert.Equal(t, suc, false)

	suc = VerifyNicknameFormat("  ")
	t.Log(suc)
	assert.Equal(t, suc, false)
}

func TestVerifyMobileRegionFormat(t *testing.T) {
	ok := VerifyMobileRegionFormat("")
	t.Log(ok)
	assert.Equal(t, ok, false)

	ok = VerifyMobileRegionFormat("+")
	t.Log(ok)
	assert.Equal(t, ok, false)

	ok = VerifyMobileRegionFormat("+86")
	t.Log(ok)
	assert.Equal(t, ok, true)

	ok = VerifyMobileRegionFormat("+8676")
	t.Log(ok)
	assert.Equal(t, ok, false)

	ok = VerifyMobileRegionFormat("86")
	t.Log(ok)
	assert.Equal(t, ok, false)
}
