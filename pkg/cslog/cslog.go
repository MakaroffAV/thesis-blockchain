package cslog

import "fmt"

func Info() {}

func Fail(err error) {
	fmt.Println(err)
}
