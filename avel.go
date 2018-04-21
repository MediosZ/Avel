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
		/*defer c.request.Body.Close()
		con, _ := ioutil.ReadAll(c.request.Body) //get data from post body
		fmt.Fprintln(c.writer, string(con))
		datafromjson := new(Somedata)
		json.Unmarshal([]byte(con), &datafromjson)
		fmt.Fprintln(c.writer, datafromjson)*/
		data := `[{"level":"debug","msg":"File Not Found","author":"Cynhard"},` +
			`{"level":"","msg":"Logic error","author":"Gopher"}]`

		var dbgInfos []DebugInfo
		json.Unmarshal([]byte(data), &dbgInfos)
		fmt.Println(dbgInfos)
	})
	mux.Post("/mhandle", func (c *Context){
		fmt.Fprintf(c.writer, "post context test, %q", html.EscapeString(c.request.URL.Path))
	})
	http.ListenAndServe(":8080", mux)
}
