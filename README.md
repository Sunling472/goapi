**Example:**

```go
package main

import (
	"encoding/json"
	"time"

	"github.com/Sunling472/goapi"
	gr "github.com/Sunling472/goapi/router"
	"github.com/Sunling472/goapi/types/methods"
)

type UrlModel struct {
	Id string
}

type QueryModel struct {
	Name string
}

type InModel struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type OutModel struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
}

func (im InModel) Json() []byte {
	data, err := json.Marshal(im)
	if err != nil {
		panic(err)
	}
	return data
}

func (om OutModel) Json() []byte {
	data, err := json.Marshal(om)
	if err != nil {
		panic(err)
	}
	return data
}

func main() {
	rt := gr.New("")
	gr.SmartHandler[methods.Get](rt, "/users/{id}",
		func(opt gr.HandlerOpts[UrlModel, QueryModel, InModel]) gr.Schema {
			return OutModel{
				Id:   opt.Url.Id,
				Name: opt.Query.Name,
			}
		},
	)

	api := goapi.New(rt, &goapi.ApiConfig{})
	api.Serve()
}
```
