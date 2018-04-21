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
}

// Form -> application/x-www-form-urlencoded && query
// PostForm -> application/x-www-form-urlencoded
//ParseMultipartForm(int) -> MultipartForm -> form-data
func (c *Context) Send(i interface{}){
	c.request.ParseForm()
	fmt.Fprintln(c.writer, c.request.Form)
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

type Handler func (c *Context)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	c := Context{
		writer: w,
		request: r,
	}
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
}

func NewMux() *Mux{
	mux := new(Mux)
	mux.ma = map[string]MethodGroup{
		"GET": make(map[string]MuxEntry),
		"POST": make(map[string]MuxEntry),
		"PUT": make(map[string]MuxEntry),
		"DELETE": make(map[string]MuxEntry),
	}
	return mux
}

func (m Mux) ServeHTTP(w http.ResponseWriter, r *http.Request){
	method := r.Method
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

func Json(i interface{}) string{
	jsondata, err := json.Marshal(i)
	if err != nil{
		return "something was wrong"
	} else{
		return string(jsondata)
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

