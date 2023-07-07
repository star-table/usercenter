package util

import "strconv"

func StrToInt64WithFallback(str string, fallback int64) int64 {
	output, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fallback
	}
	return output
}

func StrToInt64WithDefaultZero(str string) int64 {
	return StrToInt64WithFallback(str, 0)
}
