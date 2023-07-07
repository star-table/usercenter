package excel

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go/log"
	"github.com/star-table/usercenter/core/consts"
)

/**
* excelFileName 本地excel路径名
* orginFileUrl  远程excel url
* sheetIndex    第几份excel
* ignoreRows    忽略行数
* timeFormatIndex 需要转换时间格式的列
 */
func GenerateCSVFromXLSXFile(url string, urlType int, sheetIndex int, ignoreRows int, timeFormatIndex []int64) ([][]string, int, error) {
	var error error
	var xlFile *xlsx.File
	if urlType == consts.UrlTypeHttpPath {
		resp, err := http.Get(url)
		if err != nil {
			log.Error(err)
			return nil, ignoreRows, err
		}
		resBody := resp.Body
		buf := new(bytes.Buffer)
		buf.ReadFrom(resBody)

		xlFile, error = xlsx.OpenBinary(buf.Bytes())
	} else {
		xlFile, error = xlsx.OpenFile(url)

	}
	if error != nil {
		log.Error(error)
		return nil, ignoreRows, error
	}
	sheetLen := len(xlFile.Sheets)
	switch {
	case sheetLen == 0:
		return nil, ignoreRows, errors.New("This XLSX file contains no sheets.")
	case sheetIndex >= sheetLen:
		return nil, ignoreRows, fmt.Errorf("No sheet %d available, please select a sheet between 0 and %d\n", sheetIndex, sheetLen-1)
	}
	sheet := xlFile.Sheets[sheetIndex]

	result, ignoreRows := assemblyResult(sheet, ignoreRows, timeFormatIndex)

	return result, ignoreRows, nil
}

func assemblyResult(sheet *xlsx.Sheet, ignoreRows int, timeFormatIndex []int64) ([][]string, int) {

	result := [][]string{}
	if len(sheet.Rows) > 0 && len(sheet.Rows[0].Cells) > 0 {
		if len(sheet.Rows[0].Cells) == 1 || sheet.Rows[0].Cells[2].Value == "" {
			ignoreRows = 2
		} else {
			ignoreRows = 1
		}
	} else {
		ignoreRows = 1
	}
	for key, row := range sheet.Rows {
		if key <= ignoreRows-1 {
			continue
		}
		var vals []string
		if row != nil {
			for k, cell := range row.Cells {
				str, err := cell.FormattedValue()
				if err != nil {
					vals = append(vals, err.Error())
					continue
				}
				if isInIndex(timeFormatIndex, int64(k)) {
					time, err := cell.GetTime(false)
					if err != nil {
						log.Error(err)
					} else {
						str = time.Format(consts.AppTimeFormat)
					}
				}
				//vals = append(vals, fmt.Sprintf("%q", str))
				vals = append(vals, str)
			}
		}
		result = append(result, vals)
	}
	return result, ignoreRows
}

func isInIndex(param []int64, index int64) bool {
	for _, v := range param {
		if v == index {
			return true
		}
	}
	return false
}
