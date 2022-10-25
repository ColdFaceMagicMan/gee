package gee

import (
	"fmt"
	"time"
)

func Logger() HandleFunc {
	return func(c *Context) {
		t := time.Now()

		c.Next()

		fmt.Printf("[%d] %s in %v\n", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
