package asyn

import (
	"fmt"
)

func Execute(fn func()){
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("异步执行出错：%s\n", r)
			}
		}()
		fn()
	}()
}
