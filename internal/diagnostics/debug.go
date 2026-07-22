package diagnostics

import "fmt"

var Debug bool = false

func Dbg(a ...any) (int, error) {
	if Debug {
		return fmt.Println(a...)
	}
	return 0, nil
}

func Dbgfmt(format string, a ...any) (int, error) {
	if Debug {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}
