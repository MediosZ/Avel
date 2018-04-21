package main

import (
	"net/http"
	"fmt"
	"html"
)

func init() {
	fmt.Printf("this is a test server \n" +
		"listening at localhost:8080")
}

func main(){
	mux := NewMux()
	mux.Get("/", func (c *Context){
		fmt.Fprintf(c.writer, "Welcome to Avel Server, go look around if you like. \n")
	})
	mux.Get("/mhandle", func (c *Context){
		fmt.Fprintf(c.writer, "context test, %q", html.EscapeString(c.request.URL.Path))
	})
	mux.Get("/get", func (c *Context){
		c.Send("hello")
	})
	mux.Post("/mhandle", func (c *Context){
		fmt.Fprintf(c.writer, "post context test, %q", html.EscapeString(c.request.URL.Path))
	})
	http.ListenAndServe(":8080", mux)
}
