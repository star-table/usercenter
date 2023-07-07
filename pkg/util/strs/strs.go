package strs

import (
	"fmt"
	"strconv"
	"strings"
)

func Len(str string) int {
	return len([]rune(str))
}

func ObjectToString(obj interface{}) string {
	return fmt.Sprintf("%#v", obj)
}

// DesensitizeStr 将字符串的指定位置进行字符串替换。如：把手机号中间 4 位替换成 `*`。
// startIndex 位置的起始索引位置。闭区间
// endIndex 位置的结束索引位置。开区间
func DesensitizeStr(str string, startIndex, endIndex int, ele rune) string {
	rArr := []rune(str)
	newRArr := make([]rune, 0)
	for index, c := range rArr {
		if index >= startIndex && index < endIndex {
			newRArr = append(newRArr, ele)
		} else {
			newRArr = append(newRArr, c)
		}
	}
	return string(newRArr)
}

func TransferPhone(phone string, desensitizeType int) string {
	switch desensitizeType {
	case 1:
		phone = DesensitizeStr(phone, 3, 7, '*')
	default:
	}
	return phone
}

// Int64Implode 将切片内容拼接成字符串
func Int64Implode(list []int64, glue string) string {
	strList := make([]string, len(list))
	for i, item := range list {
		strList[i] = strconv.FormatInt(item, 10)
	}
	return strings.Join(strList, glue)
}