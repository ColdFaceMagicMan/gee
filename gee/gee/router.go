// 将router独立出来
package gee

import (
	"net/http"
	"strings"
)

type router struct {
	handlers map[string]HandleFunc
	roots    map[string]*node //key为get post等方法
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandleFunc),
		roots:    make(map[string]*node),
	}
}

func splitFullPath(fullPath string) []string {
	splitedPath := strings.Split(fullPath, "/")

	parts := make([]string, 0)

	for _, item := range splitedPath {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, path string, handler HandleFunc) {
	parts := splitFullPath(path)
	key := method + "-" + path

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(path, parts, 0)
	r.handlers[key] = handler
	//fmt.Printf("add %s\n", method+path)
}

// 返回node和params
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := splitFullPath(path)
	params := make(map[string]string) //gin中对param有进一步封装
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	node := root.search(searchParts, 0)

	if node != nil {
		parts := splitFullPath(node.fullPath)

		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index] // 取出动态路由片段“：xxx”在传入的path中的对应值
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}

		return node, params
	}

	return nil, nil

}

func (r *router) handle(c *Context) {

	node, params := r.getRoute(c.Method, c.Path)

	if node != nil {
		c.Params = params
		key := c.Method + "-" + node.fullPath //did a bug

		c.handlers = append(c.handlers, r.handlers[key]) //unsafe!

	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 not found")
		})

	}

	c.Next() //实际上next()方法只应该在中间件中调用
}
