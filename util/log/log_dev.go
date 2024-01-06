//go:build !prod

package log

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/gookit/color"
)

var timeColor = color.New(color.FgGreen)
var callerColor = color.C256(140)
var messageColor = color.S256(110).SetOpts(color.Opts{color.OpBold})
var argKeyColor = color.S256(246).SetOpts(color.Opts{color.OpBold})
var argValueColor = color.C256(246)
var argErrorColor = color.New(color.FgRed)
var separatorColor = color.New(color.FgGray)

func PrintTime() {
	now := time.Now().Format("03:04:05 PM")
	timeColor.Print(now)
}

func PrintCaller() {
	pc, _, line, ok := runtime.Caller(2)
	if !ok {
		panic("No caller information")
	}
	callerColor.Print(runtime.FuncForPC(pc).Name() + ":" + strconv.Itoa(line))
}

func PrintSeparator() {
	separatorColor.Print(" • ")
}

func Debug(args ...any) {
	// Print time
	PrintTime()

	// Print caller
	PrintSeparator()
	PrintCaller()

	// Print message
	if v, ok := args[0].(string); ok {
		PrintSeparator()
		messageColor.Print(v)
		args = args[1:]
	}

	// Print args
	key := true
	for len(args) > 0 {
		if v, ok := args[0].(string); ok && key {
			println()
			argKeyColor.Print(" - " + v)
			key = false
		} else if v, ok := args[0].(error); ok {
			println()
			argErrorColor.Print(" - " + v.Error())
			key = false
		} else {
			argValueColor.Print(" " + fmt.Sprintf("%v", args[0]))
			key = true
		}
		args = args[1:]
	}
	println()
}

func Error(message string, err error) {
	// Print time
	PrintTime()

	// Print caller
	PrintSeparator()
	PrintCaller()

	// Print message
	PrintSeparator()
	messageColor.Print(message)

	// Print error
	PrintSeparator()
	argErrorColor.Print(err.Error())
}
