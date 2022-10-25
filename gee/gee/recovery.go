package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func trace(message string) string {
	var pcs [32]uintptr

	//跳过前三个（caller、trace、recovery）
	n := runtime.Callers(3, pcs[:])

	//A Builder is used to efficiently build a string using Write methods.
	//It minimizes memory copying. The zero value is ready to use. Do not copy a non-zero Builder.
	var str strings.Builder
	str.WriteString(message + "\ntraceback:")
	for _, pc := range pcs[:n] {
		//FuncForPC returns a *Func describing the function that contains the given program counter address, or else nil.
		fn := runtime.FuncForPC(pc)
		//FileLine returns the file name and line number of the source code corresponding to the program counter pc.
		//The result will not be accurate if pc is not a program counter within f.
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

// 中间件，调用完next后检查错误
func Recovery() HandleFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal server error")
			}
		}()

		c.Next()
	}
}
