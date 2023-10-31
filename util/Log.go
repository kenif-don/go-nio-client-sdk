package util

import (
	"fmt"
)

// Out 统一的日志打印函数
func Out(str string, params ...interface{}) {
	fmt.Printf(str, params)
}

// Err 统一的错误日志打印函数
func Err(str string, params ...interface{}) {
	fmt.Errorf(str, params)
}
