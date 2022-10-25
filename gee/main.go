package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()

	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello world\n")
	})

	r.GET("/panic", func(c *gee.Context) {
		names := []string{}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":8000")
}
