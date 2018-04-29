# Avel

Avel is a light server framework based on Go, aiming to build backend server quickly.
Inspired by martini https://github.com/go-martini/martini



## Usage

### Basic usage

```go
// create the application
Mux := NewMux()
// add route
Mux.Get("/", func (c *Context){
    c.Send("Hello, world!")
})
// Start Server, listen at localhost:8080 by default
Mux.Listen()
```

### Methods

Get, Post, Put, Delete are supported.

### Middleware

The func will be called in the order you defined.

The Next() function will just call the next middleware and then go back.

If you use Use(), the middleware will be called first.

```Go
Mux.Use(middlerware func (c *Context))

Mux.Use(func (c *Context){
	c.Send("hello, this is a middleware")
	c.Next()
	c.Send("after next func")
})

Mux.Use(func (c *Context){
	c.Send("hello, this is a middleware2")
	c.Next()
	c.Send("after next2")
})

Mux.Use(func (c *Context){
	c.Send("hello, this is a middleware3")
	c.Next()
	c.Send("after next3")
})

// you will get following result
hello, this is a middleware
hello, this is a middleware2
hello, this is a middleware3
after next3
after next2
after next func
```

Also, you can just have middleware worked on a single path.

```Go
mux.Get("/mid", func1, func2, func3)
```

### Data format

You can send json using Json() and decode json data with Decode()



## Todo-List

Depends on the need.