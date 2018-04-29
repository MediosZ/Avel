package avel

import (
	"net/http"
	"sync"
	"fmt"
	"html"
	"encoding/json"
)

type Context struct {
	writer http.ResponseWriter
	request *http.Request
	handlers []Hand
	index int
}

// Form -> application/x-www-form-urlencoded && query
// PostForm -> application/x-www-form-urlencoded
//ParseMultipartForm(int) -> MultipartForm -> form-data
func (c *Context) Send(i interface{}){
	c.request.ParseForm()
	//type switch
	switch v := i.(type) {
	case int:
		fmt.Fprintln(c.writer, v)
	case string:
		fmt.Fprintln(c.writer, v)
	default:{
		if jsondata,err := json.Marshal(i); err != nil{
			fmt.Fprintln(c.writer, "something was wrong")
		} else {
			fmt.Fprintln(c.writer, string(jsondata))
		}}
	}
}

//这里传递的context其实是Hand.ctx

func (c *Context) Next(){
	if((c.index + 1) < len(c.handlers)) {
		c.index += 1
		c.handlers[c.index].ServeHTTP(c.writer, c.request)
	}else {
		return
	}
}

func (c *Context) run(){
	c.index = 0
	for c.index < len(c.handlers){
		c.handlers[c.index].ServeHTTP(c.writer, c.request)
		c.index += 1
	}
}


type Hand struct{
	h Handler
	ctx *Context
}

func (h *Hand) ServeHTTP(w http.ResponseWriter, r *http.Request){
	h.ctx.writer = w
	h.ctx.request = r
	defer r.Body.Close()
	h.h(h.ctx)
}

type Handler func (c *Context)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	c := Context{
		writer: w,
		request: r,
	}
	defer r.Body.Close()
	h(&c)
}


type MuxEntry struct{
	method string
	pat string
	ctx *Context

}

func (me *MuxEntry) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	me.ctx.writer = w
	me.ctx.request = r
	me.ctx.run()
}

type MethodGroup map[string]MuxEntry

// Mux is main Route
type Mux struct{
	mu sync.RWMutex
	ma map[string]MethodGroup
	ctx *Context
}

func NewMux() *Mux{
	mux := new(Mux)
	mux.ma = map[string]MethodGroup{
		"GET": make(map[string]MuxEntry),
		"POST": make(map[string]MuxEntry),
		"PUT": make(map[string]MuxEntry),
		"DELETE": make(map[string]MuxEntry),
	}
	mux.ctx = &Context{
		handlers: make([]Hand, 0),
		index: 0,
	}
	return mux
}

func (m *Mux) Use(fn func(c *Context)){
	m.ctx.handlers = append(m.ctx.handlers, Hand{
		h: Handler(fn),
		ctx: m.ctx,
	})
}

func (m *Mux) UseAll(fnarr []func(c *Context)){
	for _, fn := range fnarr{
		m.Use(fn)
	}
}


func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request){
	m.ctx.writer = w
	m.ctx.request = r
	method := r.Method

	m.ctx.run()

	if me, ok := m.ma[method][r.URL.Path]; !ok {
		mh := Handler(func(c *Context){
			fmt.Fprintf(c.writer, "Hello, but there is no handler for %q \n", html.EscapeString(c.request.URL.Path))
			fmt.Fprintln(c.writer, "try another or just read the doc")
		})
		mh.ServeHTTP(w, r)
	} else{
		me.ServeHTTP(w, r)
	}

}


func (mu *Mux) Get(pattern string, fns ...func(c *Context)){
	mue := MuxEntry{
		method: "GET",
		pat: pattern,
		ctx: &Context{
			handlers: make([]Hand, 0),
			index: 0,
		},
	}
	for _, fn := range fns{
		mue.ctx.handlers = append(mue.ctx.handlers, Hand{
			h: Handler(fn),
			ctx: mue.ctx,
		})
	}

	mu.ma["GET"][pattern] = mue
}

func (mu *Mux) Post(pattern string, fns ...func(c *Context)){
	mue := MuxEntry{
		method: "POST",
		pat: pattern,
		ctx: &Context{
			handlers: make([]Hand, 0),
			index: 0,
		},
	}
	for _, fn := range fns{
		mue.ctx.handlers = append(mue.ctx.handlers, Hand{
			h: Handler(fn),
			ctx: mue.ctx,
		})
	}
	mu.ma["POST"][pattern] = mue
}

func (mu *Mux) Put(pattern string, fns ...func(c *Context)){
	mue := MuxEntry{
		method: "PUT",
		pat: pattern,
		ctx: &Context{
			handlers: make([]Hand, 0),
			index: 0,
		},
	}
	for _, fn := range fns{
		mue.ctx.handlers = append(mue.ctx.handlers, Hand{
			h: Handler(fn),
			ctx: mue.ctx,
		})
	}
	mu.ma["Put"][pattern] = mue
}

func (mu *Mux) Delete(pattern string, fns ...func(c *Context)){
	mue := MuxEntry{
		method: "DELETE",
		pat: pattern,
		ctx: &Context{
			handlers: make([]Hand, 0),
			index: 0,
		},
	}
	for _, fn := range fns{
		mue.ctx.handlers = append(mue.ctx.handlers, Hand{
			h: Handler(fn),
			ctx: mue.ctx,
		})
	}
	mu.ma["DELETE"][pattern] = mue
}

func (m *Mux) ListenAt(port string){
	http.ListenAndServe(port, m)
}

func (m *Mux) Listen(){
	http.ListenAndServe(":8080", m)
}

