package slice

import (
	"errors"
	"reflect"
	"sort"
)

//通过map去重slice
func SliceUniqueString(s []string) []string {
	res := make([]string, 0)
	exist := make(map[string]bool)
	for _, s2 := range s {
		if _, ok := exist[s2]; ok {
			continue
		}

		res = append(res, s2)
		exist[s2] = true
	}

	return res
}

func SliceUniqueInt64(s []int64) []int64 {
	res := make([]int64, 0)
	exist := make(map[int64]bool)
	for _, i2 := range s {
		if _, ok := exist[i2]; ok {
			continue
		}
		res = append(res, i2)
		exist[i2] = true
	}

	return res
}

func Contain(list interface{}, obj interface{}) (bool, error) {
	targetValue := reflect.ValueOf(list)
	switch reflect.TypeOf(list).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
		return false, nil
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
		return false, nil
	}
	return false, errors.New("not in array")
}

func ToSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		panic("toslice arr not slice")
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}

func SearchInt64(a []int64, x int64) int {
	return sort.Search(len(a), func(i int) bool { return a[i] == x })
}

// Int64Intersect 两个 int64 数组切片的交集
func Int64Intersect(list1, list2 []int64) []int64 {
	uniqueMap := map[int64]bool{}
	for _, ele := range list1 {
		uniqueMap[ele] = true
	}
	result := make([]int64, 0)
	for _, ele := range list2 {
		if _, ok := uniqueMap[ele]; ok {
			result = append(result, ele)
		}
	}
	return result
}

// 查找在 arr1 中，但不在 arr2 中的元素
func ArrayDiffInt64(arr1, arr2 []int64) (diffArr []int64) {
	if len(arr2) < 1 || len(arr1) < 1 {
		diffArr = arr1
		return
	}
	for i := 0; i < len(arr1); i++ {
		item := arr1[i]
		isIn := false
		for j := 0; j < len(arr2); j++ {
			if item == arr2[j] {
				isIn = true
				break
			}
		}
		if !isIn {
			diffArr = append(diffArr, item)
		}
	}
	return diffArr
}

// ArrayDiffString 查找在 arr1 中，但不在 arr2 中的元素
func ArrayDiffString(arr1, arr2 []string) (diffArr []string) {
	if len(arr2) < 1 || len(arr1) < 1 {
		diffArr = arr1
		return
	}
	for i := 0; i < len(arr1); i++ {
		item := arr1[i]
		isIn := false
		for j := 0; j < len(arr2); j++ {
			if item == arr2[j] {
				isIn = true
				break
			}
		}
		if !isIn {
			diffArr = append(diffArr, item)
		}
	}
	return diffArr
}

// 删除切片中的元素
func DeleteOneElemOfSlice(arr []string, elem string) []string {
	tmp := make([]string, 0, len(arr))
	for _, v := range arr {
		if v != elem {
			tmp = append(tmp, v)
		}
	}
	return tmp
}

// 判断数字切片是否有某一元素
func IntContains(elems []int64, v int64) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
