package date

import (
	"fmt"
	"testing"
	"time"

	"github.com/star-table/usercenter/core/types"
)

func TestFormatTime(t *testing.T) {

	t1 := types.Time(time.Now())

	str := FormatTime(t1)
	fmt.Println(str)
	fmt.Println(ParseTime(str))

}

func TestFormat(t *testing.T) {

}
