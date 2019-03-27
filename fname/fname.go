package fname

import (
	"runtime"
	"strings"
)

func Current() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0]).Name()

	spName := strings.Split(f, "/")
	if len(spName) > 1 {
		return spName[len(spName)-1]
	}

	return f
}
