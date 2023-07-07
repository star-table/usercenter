package random

import "testing"

func TestGenerateRandomStringAsBase64(t *testing.T) {
	t.Log(GenerateRandomStringAsBase36(32))
}
