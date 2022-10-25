// 实现中间间middleware
package gee

import (
	"html/template"
	"net/http"
	"strings"
)

type HandleFunc func(*Context)

// 内核基础为路径到HandleFunc的映射
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup

	//for html render
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 添加route , 没有处理重复监听
func (e *Engine) addRoute(method string, pattern string, handler HandleFunc) {
	e.router.addRoute(method, pattern, handler)
}

func (e *Engine) GET(pattern string, handler HandleFunc) {
	e.addRoute("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandleFunc) {
	e.addRoute("POST", pattern, handler)
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

func (group *RouterGroup) Use(middleWares ...HandleFunc) {
	group.middleWares = append(group.middleWares, middleWares...)
}

// ListenAndServe将请求交给接口的ServeHTTP处理
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middleWares []HandleFunc

	//urlpath含有该组前缀时将该组的中间件加入函数链
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middleWares = append(middleWares, group.middleWares...)
		}
	}

	c := NewContext(w, req)
	c.handlers = middleWares
	c.engine = e
	e.router.handle(c)
}

// http.FileSystem为可以open的interface
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandleFunc {
	absolutePath := group.prefix + relativePath

	//跳过前缀，留下文件路径交给FileServer返回的handler
	//To use the operating system's file system implementation, use http.Dir:
	//http.Handle("/", http.FileServer(http.Dir("/tmp")))
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.GetParam("filePath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	fileHandler := group.createStaticHandler(relativePath, http.Dir(root))
	pattern := relativePath + "/*filePath"
	group.GET(pattern, fileHandler)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}
