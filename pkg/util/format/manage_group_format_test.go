package format

import (
	"gotest.tools/assert"
	"testing"
)

func TestVerifyAndGetAppOrPkgIdFormat(t *testing.T) {

	tag, id, ok := VerifyAndGetAppOrPkgIdFormat("1234")
	t.Log(tag, id, ok)
	assert.Equal(t, tag, "")
	assert.Equal(t, id, int64(0))
	assert.Equal(t, ok, false)

	tag, id, ok = VerifyAndGetAppOrPkgIdFormat("p1234p")
	t.Log(tag, id, ok)
	assert.Equal(t, tag, "")
	assert.Equal(t, id, int64(0))
	assert.Equal(t, ok, false)

	tag, id, ok = VerifyAndGetAppOrPkgIdFormat("ap1234")
	t.Log(tag, id, ok)
	assert.Equal(t, tag, "")
	assert.Equal(t, id, int64(0))
	assert.Equal(t, ok, false)

	tag, id, ok = VerifyAndGetAppOrPkgIdFormat("p123412412a3")
	t.Log(tag, id, ok)
	assert.Equal(t, tag, "")
	assert.Equal(t, id, int64(0))
	assert.Equal(t, ok, false)

	tag, id, ok = VerifyAndGetAppOrPkgIdFormat("p1234124123")
	t.Log(tag, id, ok)
	assert.Equal(t, tag, "p")
	assert.Equal(t, id, int64(1234124123))
	assert.Equal(t, ok, true)

	tag, id, ok = VerifyAndGetAppOrPkgIdFormat("a1234124123")
	t.Log(tag, id, ok)
	assert.Equal(t, tag, "a")
	assert.Equal(t, id, int64(1234124123))
	assert.Equal(t, ok, true)
}
