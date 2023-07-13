package util

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func GetCallers() string {
	callers := ""
	for i := 0; true; i++ {
		_, file, line, ok := runtime.Caller(i)
		file = filepath.Base(file)
		if !ok {
			break
		}
		callers = callers + fmt.Sprintf("cause: [%v:%v]\n", file, line)
	}
	return callers
}
