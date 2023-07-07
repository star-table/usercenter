package str

import (
	"encoding/json"
	"strconv"

	"github.com/star-table/usercenter/pkg/util/strs"

	"regexp"
	"strings"
	"unicode"

	"github.com/axgle/mahonia"
)

func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

//start：正数 - 在字符串的指定位置开始,超出字符串长度强制把start变为字符串长度
//       负数 - 在从字符串结尾的指定位置开始
//       0 - 在字符串中的第一个字符处开始
//length:正数 - 从 start 参数所在的位置返回
//       负数 - 从字符串末端返回
func Substr(str string, start, length int) string {
	if length == 0 {
		return ""
	}
	runeStr := []rune(str)
	len_str := len(runeStr)

	if start < 0 {
		start = len_str + start
	}
	if start > len_str {
		start = len_str
	}
	end := start + length
	if end > len_str {
		end = len_str
	}
	if length < 0 {
		end = len_str + length
	}
	if start > end {
		start, end = end, start
	}
	return string(runeStr[start:end])
}

func UrlParse(url string) (string, string) {
	var host, path string
	split := strings.Split(url, "/")
	if split[0] != "http:" && split[0] != "https:" {
		return host, url
	}
	host = strings.Join(split[:3], "/")
	path = "/" + strings.Join(split[3:], "/")
	return host, path
}

func ParseOssKey(url string) string {
	_, path := UrlParse(url)
	if strings.HasPrefix(path, "/") {
		return path[1:strs.Len(path)]
	}
	return path
}

func AppendJSONSliceStr(str, jsonStr string) string {

	//只能判断已经有值了判断不了是slice还是string
	traceIds := &[]string{}
	//解析判断是否是数组
	err := json.Unmarshal([]byte(jsonStr), traceIds)
	if err == nil {
		//解析成功 原本已经是数组了,拼到后面去然后在抓json
		*traceIds = append(*traceIds, str)
		newTraceIds, _ := json.MarshalIndent(traceIds, "", "	")
		return string(newTraceIds)
	}

	//解析不成功,走到这边说明原本是个string 拼到数组里面去
	*traceIds = append(*traceIds, str, jsonStr)
	newTraceIds, _ := json.MarshalIndent(traceIds, "", "	")
	return string(newTraceIds)
}

func CountStrByGBK(str string) int {
	new := mahonia.NewEncoder("gbk").ConvertString(str)
	return len(new)
}

func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//将p标签换成换行符
	re, _ = regexp.Compile("\\<p[^>]*>")
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile("</p>")
	src = re.ReplaceAllString(src, "\n")
	// 将br标签换成换行符
	re, _ = regexp.Compile("<br/>")
	src = re.ReplaceAllString(src, "\n")
	// 替换img为（url）
	re, _ = regexp.Compile("\\<img[^>]*src=(\\'|\")(.*?)[^>]*>")
	for _, s := range re.FindAllString(src, -1) {
		re, _ = regexp.Compile("src=(\\'|\")([^\\'\"]*)(\\'|\")")
		url := re.FindAllString(s, -1)
		if len(url) == 1 {
			src = strings.ReplaceAll(src, s, "[附件]("+url[0][5:len(url[0])-1]+")")
		}
	}
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "")
	//去除连续的换行符
	//re, _ = regexp.Compile("\\s{2,}")
	//src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

// Int64Implode 将切片内容拼接成字符串
func Int64Implode(list []int64, glue string) string {
	strList := make([]string, len(list))
	for i, item := range list {
		strList[i] = strconv.FormatInt(item, 10)
	}
	return strings.Join(strList, glue)
}
