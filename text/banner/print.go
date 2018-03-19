package banner

import "fmt"

var defaultFont *Font

func init() {
	defaultFont, _ = NewFontString(fontRectangles)
}

func Printf(format string, a ...interface{}) {
	res := fmt.Sprintf(format, a...)
	fmt.Print(New(defaultFont, res).String())
}

func Sprintf(format string, a ...interface{}) string {
	res := fmt.Sprintf(format, a...)
	return fmt.Sprint(New(defaultFont, res).String())
}
