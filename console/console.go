package console

import (
	"fmt"
	"github.com/mgutz/ansi"
	"os"
	"strings"
)

// colorOut 内部使用，设置高亮颜色
func colorOut(message, color string) {
	fmt.Fprintln(os.Stdout, ansi.Color(message, color))
}

// Success 打印一条成功消息，绿色输出
func Success(msg string) {
	colorOut(msg, "green")
}

// SuccessPretty 打印一条成功消息，绿色输出
func SuccessPretty(title, msg string) {
	tLen := len(title)
	mLen := len(msg)
	maxLen := tLen

	if maxLen < mLen {
		maxLen = mLen
	}

	colorOut(strings.Repeat("-", maxLen), "green")

	if maxLen == tLen {
		colorOut(title, "green")
	} else {
		colorOut(strings.Repeat(" ", (maxLen-tLen)/2)+title+strings.Repeat(" ", (maxLen-tLen)/2), "green")
	}

	if maxLen == mLen {
		colorOut(msg, "green")
	} else {
		colorOut(strings.Repeat(" ", (maxLen-mLen)/2)+msg+strings.Repeat(" ", (maxLen-mLen)/2), "green")
	}

	colorOut(strings.Repeat("-", maxLen), "green")
}

// Error 打印一条报错消息，红色输出
func Error(msg string) {
	colorOut(msg, "red")
}

// Warning 打印一条提示消息，黄色输出
func Warning(msg string) {
	colorOut(msg, "yellow")
}

// Exit 打印一条报错消息，并退出 os.Exit(1)
func Exit(msg string) {
	Error(msg)
	os.Exit(1)
}

// ExitIf 语法糖，自带 err != nil 判断
func ExitIf(err error) {
	if err != nil {
		Exit(err.Error())
	}
}
