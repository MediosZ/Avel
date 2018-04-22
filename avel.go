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
		defer c.request.Body.Close()

		//respbody := make(map[string]interface{})
		sdata := new(Somedata)
		//对返回的结果response进行json解码，可以利用map[string]interface{} 通用的格式存放解码后的值，或者使用具体的response的返回的信息，构建一个结构体，来存放解析后的返回值
		json.NewDecoder(c.request.Body).Decode(&sdata)
		fmt.Fprintf(c.writer, "%+v", sdata.Data)

	})
	mux.Post("/mhandle", func (c *Context){
		fmt.Fprintf(c.writer, "post context test, %q", html.EscapeString(c.request.URL.Path))
	})
	http.ListenAndServe(":8080", mux)
}
