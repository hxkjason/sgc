package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// GetCommonAddDeleteUintItems []uint 获取新旧两数组中的相同、新增、删除的元素集合
func GetCommonAddDeleteUintItems(arrNew, arrOld []uint) (commonItems, addItems, deleteItems []uint) {
	arrNew = RemoveDuplicationUint(arrNew)
	arrOld = RemoveDuplicationUint(arrOld)

	commonItems = GetSameItemsFromUintArr(arrNew, arrOld)
	if len(commonItems) == 0 {
		addItems = arrNew
		deleteItems = arrOld
	} else {
		addItems = DiffUintArrAttrs(arrNew, commonItems)
		deleteItems = DiffUintArrAttrs(arrOld, commonItems)
	}
	return commonItems, addItems, deleteItems
}

// GetSameItemsFromUintArr 获取两整型数组公有元素组成的数组
func GetSameItemsFromUintArr(arr1, arr2 []uint) []uint {
	var intersection []uint
	arr3 := append(arr1, arr2...)
	sameElem := make(map[uint]struct{})

	for _, v := range arr3 {
		if _, ok := sameElem[v]; ok {
			intersection = append(intersection, v)
		} else {
			sameElem[v] = struct{}{}
		}
	}
	return intersection
}

// []uint 获取第一个数组在第二个数组中没有的元素 即第一个数组比第二个数组多哪些

func DiffUintArrAttrs(arr1, arr2 []uint) []uint {

	var diffArr []uint
	for _, item1 := range arr1 {
		var inArray2 bool
		for _, item2 := range arr2 {
			if item1 == item2 {
				inArray2 = true
				break
			}
		}
		if !inArray2 {
			diffArr = append(diffArr, item1)
		}
	}
	return diffArr
}

// 获取两数组公有元素组成的数组

func SameArrAttrs(arr1, arr2 []string) []string {
	var intersection []string
	arr3 := append(arr1, arr2...)
	sameElem := make(map[string]struct{})

	for _, v := range arr3 {
		if _, ok := sameElem[v]; ok {
			intersection = append(intersection, v)
		} else {
			sameElem[v] = struct{}{}
		}
	}
	return intersection
}

// 获取第一个数组在第二个数组中没有的元素

func DiffArrAttrs(arr1, arr2 []string) []string {

	// 数组2 => map
	var arr2Map = make(map[string]struct{}, len(arr2))
	for _, u := range arr2 {
		arr2Map[u] = struct{}{}
	}

	var arr1OnlyHasItems = make([]string, 0, len(arr1))
	for _, item := range arr1 {
		if _, ok := arr2Map[item]; !ok {
			arr1OnlyHasItems = append(arr1OnlyHasItems, item)
		}
	}
	return arr1OnlyHasItems
}

// 获取第一个数组在第二个数组中没有的元素  第一个数组比第二个数组多哪些元素

func FirstArrOnlyHas(arr1, arr2 []uint) []uint {

	// 数组1 这里不做去重处理
	// 数组2 => map
	var arr2Map = make(map[uint]struct{}, len(arr2))
	for _, u := range arr2 {
		arr2Map[u] = struct{}{}
	}

	var arr1OnlyHasItems = make([]uint, 0, len(arr1))
	for _, item := range arr1 {
		if _, ok := arr2Map[item]; !ok {
			arr1OnlyHasItems = append(arr1OnlyHasItems, item)
		}
	}
	return arr1OnlyHasItems
}

// 第一个数组元素全在第二个数组范围内

func FirstArrAllInSecondArr(arr1, arr2 []uint) bool {
	// 数组2 => map
	var arr2Map = make(map[uint]struct{}, len(arr2))
	for _, u := range arr2 {
		arr2Map[u] = struct{}{}
	}

	for _, item := range arr1 {
		if _, ok := arr2Map[item]; !ok {
			return false
		}
	}
	return true
}

// RemoveDuplicationString 去除重复元素 string
func RemoveDuplicationString(arr []string) []string {
	set := make(map[string]struct{}, len(arr))
	j := 0
	for _, v := range arr {
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = struct{}{}
		arr[j] = v
		j++
	}

	return arr[:j]
}

// RemoveDuplicationUint 去除重复元素 uint
func RemoveDuplicationUint(arr []uint) []uint {
	set := make(map[uint]struct{}, len(arr))
	j := 0
	for _, v := range arr {
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = struct{}{}
		arr[j] = v
		j++
	}

	return arr[:j]
}

// UintArrToBlockArr  将整型数组切块
func UintArrToBlockArr(arr []uint, blockNum int) [][]uint {

	var result [][]uint
	arrLen := len(arr)

	if blockNum >= arrLen {
		result = append(result, arr)
		return [][]uint{arr}
	}

	startIndex := 0
	endIndex := blockNum

	for {
		result = append(result, arr[startIndex:endIndex])
		if endIndex == arrLen {
			break
		}
		startIndex += blockNum
		endIndex = startIndex + blockNum
		if endIndex >= arrLen {
			endIndex = arrLen
		}
	}

	return result
}

// StringArrToBlockArr  将字符串数组切块
func StringArrToBlockArr(arr []string, blockNum int) [][]string {

	var result [][]string
	arrLen := len(arr)

	if blockNum >= arrLen {
		return [][]string{arr}
	}

	startIndex := 0
	endIndex := blockNum

	for {
		result = append(result, arr[startIndex:endIndex])
		if endIndex == arrLen {
			break
		}
		startIndex += blockNum
		endIndex = startIndex + blockNum
		if endIndex >= arrLen {
			endIndex = arrLen
		}
	}

	return result
}

// InArray 判断元素是否在对象中 object supported types: slice, array or map
func InArray(item interface{}, object interface{}) bool {
	val := reflect.ValueOf(object)
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(item, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(item, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		fmt.Printf("object: {%v} must be slice, array or map ===============\n", object)
		return false
	}
	return false
}

// CheckHasRepeatItem 判断数组或切片是否有重复元素
func CheckHasRepeatItem(object interface{}) (hasRepeat bool, repeatI string) {
	val := reflect.ValueOf(object)
	repeat := false
	repeatItem := ""
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		sliceLen := val.Len()
		for i := 0; i < sliceLen; i++ {
			for j := i + 1; j < sliceLen; j++ {
				if val.Index(i).Interface() == val.Index(j).Interface() {
					repeat = true
					repeatItem = val.Index(j).String()
					break
				}
			}
			if repeat {
				break
			}
		}
	default:
		panic(errors.New("object type must be slice or array"))
	}
	return repeat, repeatItem
}

// UintSliceToString []uint{1,2,3} => 1,2,3
func UintSliceToString(uintSlice []uint) string {

	if len(uintSlice) == 0 {
		return ""
	}

	b := make([]string, len(uintSlice))
	for i, v := range uintSlice {
		b[i] = strconv.Itoa(int(v))
	}
	return strings.Join(b, ",")
	//return strings.Trim(strings.Replace(fmt.Sprint(a), " ", ",", -1), "[]")
}
