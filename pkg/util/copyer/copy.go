package copyer

import (
	"github.com/pkg/errors"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/pkg/util/json"
)

func Copy(src interface{}, source interface{}) error {
	jsonStr, err := json.ToJson(src)
	if err != nil {
		return errors.New("json转换异常")
	}
	err = json.FromJson(jsonStr, source)
	if err != nil {
		return errs.JSONConvertError
	}
	return nil
}
