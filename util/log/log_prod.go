//go:build prod

package log

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

func Debug(args ...any) {
}

func Error(message string, err error) {
	now := time.Now().Format("03:04:05 PM")
	fmt.Print("[" + now)
	pc, _, line, ok := runtime.Caller(1)
	if !ok {
		panic("No caller information")
	}
	fmt.Print(" " + runtime.FuncForPC(pc).Name() + ":" + strconv.Itoa(line) + "]")
	fmt.Print(" " + message + " | " + err.Error())
}
