package excel

import (
	"fmt"
	"testing"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/pkg/util/slice"
)

func TestGenerateCSVFromXLSXFile(t *testing.T) {
	t.Log(GenerateCSVFromXLSXFile("issue.xlsx", consts.UrlTypeDistPath, 0, 0, []int64{1, 2}))
}

func TestGenerateCSVFromXLSXFile2(t *testing.T) {
	a := []int64{1, 2}
	fmt.Println(slice.Contain(a, 1))
}
