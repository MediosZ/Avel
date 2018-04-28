package main

import (
	"net/http"
	"fmt"
	"html"
	"encoding/json"
)

func init() {
	fmt.Printf("this is a test server \n" +
		"listening at localhost:8080\n")
}

type Somedata struct{
	Data string `json:"data"`
	Desp string `json:"desp"`
}

type DebugInfo struct {
	Level  string `json:"level,omitempty"` // Level解析为level,忽略空值
	Msg    string `json:"message"`         // Msg解析为message
	Author string `json:"author"`               // 忽略Author
}


func main(){
	deb := DebugInfo{
		Level: "top",
		Msg: "hello there",
		Author: "tricster",
	}
	sdd := Somedata{Data: "hello world",Desp: "this is a test data"}
	jsondata, _ := json.Marshal(deb)
	mux := NewMux()

	mux.Use(func (c *Context){
		c.Send("hello, this is a middleware")
		c.Next()
		c.Send("after next func")
	})

	mux.Use(func (c *Context){
		c.Send("hello, this is a middleware2")
		c.Next()
		c.Send("after next2")
	})

	mux.Use(func (c *Context){
		c.Send("hello, this is a middleware3")
		c.Next()
		c.Send("after next3")
	})

	mux.Get("/", func (c *Context){
		fmt.Fprintf(c.writer, "Welcome to Avel Server, go look around if you like. \n")
	})
	mux.Get("/mhandle", func (c *Context){
		fmt.Fprintf(c.writer, "context test, %q", html.EscapeString(c.request.URL.Path))
	})
	mux.Get("/get", func (c *Context){
		c.Send(Json(sdd))
	})
	mux.Get("/num", func (c *Context){
		c.Send(deb)
	})
	mux.Post("/num", func (c *Context){
		c.Send(1234)
	})
	mux.Get("/json", func (c *Context){
		c.Send(string(jsondata))
	})
	mux.Post("/json", func (c *Context){
		respbody := Decode(c.request.Body)
		fmt.Fprintf(c.writer, "%+v", respbody)
	})
	mux.Post("/mhandle", func (c *Context){
		fmt.Fprintf(c.writer, "post context test, %q", html.EscapeString(c.request.URL.Path))
	})
	http.ListenAndServe(":8080", mux)
}
