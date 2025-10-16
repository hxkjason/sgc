package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unsafe"
)

var (
	NumericReg = regexp.MustCompile("[0-9]+")
)

// Md5 md5
func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func StrToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// SplicingStr 拼接字符串
func SplicingStr(str ...string) string {
	var build strings.Builder
	for _, v := range str {
		build.WriteString(v)
	}
	return build.String()
}

// GetNumericStr 从字符串中提取数字
func GetNumericStr(str string) string {
	return strings.Join(NumericReg.FindAllString(str, -1), "")
}

// StructToJsonBytes struct to json bytes
func StructToJsonBytes(structContent interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(structContent)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func DebugStruct(name string, structContent interface{}) {
	fmt.Println("=====" + name + "==================== start")
	jsonBytes, _ := StructToJsonBytes(structContent)
	fmt.Println(BytesToStr(jsonBytes))
	fmt.Println("=================================== end")
}

func PrintMarshalIndent(structContent interface{}, name ...interface{}) {
	fmt.Println("================================== start")
	for _, n := range name {
		fmt.Println(n)
	}
	jsonBytes, err := json.MarshalIndent(structContent, "", " ")
	if err != nil {
		fmt.Println("hasErr:", err.Error())
	} else {
		fmt.Println(BytesToStr(jsonBytes))
	}
	fmt.Println("=================================== end")
}
