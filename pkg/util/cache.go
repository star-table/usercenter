package util

import (
	"github.com/star-table/usercenter/pkg/util/temp"
)

func ParseCacheKey(key string, params map[string]interface{}) (string, error) {
	target, err := temp.Render(key, params)
	if err != nil {
		return "", err
	}
	return target, nil
}
