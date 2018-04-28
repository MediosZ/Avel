package main

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
	//fmt.Fprintln(c.writer, c.request.Form)
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

type Handler func (c *Context)

func (h Hand) ServeHTTP(w http.ResponseWriter, r *http.Request){
	h.ctx.writer = w
	h.ctx.request = r
	defer r.Body.Close()
	h.h(h.ctx)
}

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
	h http.Handler
}

type MethodGroup map[string]MuxEntry

type Mux struct{
	mu sync.RWMutex
	ma map[string]MethodGroup
	ctx *Context
}

func (m *Mux) Use(fn func(c *Context)){
	m.ctx.handlers = append(m.ctx.handlers, Hand{
		h: Handler(fn),
		ctx: m.ctx,
	})
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
		me.h.ServeHTTP(w, r)
	}

}



func (mu *Mux) Get(pattern string, fn func(c *Context)){
	mue := MuxEntry{
		method: "GET",
		pat: pattern,
		h: Handler(fn),
	}
	mu.ma["GET"][pattern] = mue
}

func (mu *Mux) Post(pattern string, fn func(c *Context)){
	mue := MuxEntry{
		method: "POST",
		pat: pattern,
		h: Handler(fn),
	}
	mu.ma["POST"][pattern] = mue
}

func (mu *Mux) Put(pattern string, fn func(c *Context)){
	mue := MuxEntry{
		method: "PUT",
		pat: pattern,
		h: Handler(fn),
	}
	mu.ma["Put"][pattern] = mue
}

func (mu *Mux) Delete(pattern string, fn func(c *Context)){
	mue := MuxEntry{
		method: "DELETE",
		pat: pattern,
		h: Handler(fn),
	}
	mu.ma["DELETE"][pattern] = mue
}

