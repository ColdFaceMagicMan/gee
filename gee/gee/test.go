package gee

import (
	"fmt"
	"log"
	"reflect"
)

func TestParsePattern() {
	ok := reflect.DeepEqual(splitFullPath("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(splitFullPath("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(splitFullPath("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		log.Fatal("test parsePattern failed")
	}
}

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestGetRoute() {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/")

	if n == nil {
		log.Fatal("nil shouldn't be returned")
	}

	if n.fullPath != "/assets/*filepath" {
		log.Fatal("should match /hello/:name")
	}

	if ps["filepath"] != "geektutu/123" {
		log.Fatal("name should be equal to 'geektutu'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.fullPath, ps["name"])

}
