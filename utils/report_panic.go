package utils

import "fmt"

// ReportPanic recovers from a panic and prints the underlying error if one occurred.
func ReportPanic() {
	p := recover()
	if p == nil {
		return
	}
	err, ok := p.(error)
	if ok {
		fmt.Println(err)
	} else {
		panic(p)
	}
}
