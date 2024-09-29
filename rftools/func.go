package rftools

import (
	"encoding/json"
	"fmt"
	"net/http"
	rf "reflect"
	"strconv"
)

type ArgsIn struct {
	Name string
	Type rf.Type
	Kind rf.Kind
}

type ArgsOut struct {
	Type rf.Type
	Kind rf.Kind
}

type Signature struct {
	In  []ArgsIn
	Out []ArgsOut
}

func GetSignature(f any) Signature {
	fType := rf.TypeOf(f)

	if fType.Kind() != rf.Func {
		panic("arg allow only function")
	}

	argsIn := []ArgsIn{}
	if fType.NumIn() > 0 {
		for i := 0; i < fType.NumIn(); i++ {
			arg := fType.In(i)
			argsIn = append(argsIn, ArgsIn{Type: arg, Name: arg.Name(), Kind: arg.Kind()})
		}
	}

	argOut := []ArgsOut{}
	if fType.NumOut() > 0 {
		for i := 0; i < fType.NumOut(); i++ {
			arg := fType.Out(i)
			argOut = append(argOut, ArgsOut{Type: arg, Kind: arg.Kind()})
		}
	}

	return Signature{
		In:  argsIn,
		Out: argOut,
	}
}

type IHandlerOpts interface{}
type OutSchema interface {
	toJson() []byte
}

type Pat struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

var Hdls []Pat

func thandle[T any](method string, path string, h func(opt T) OutSchema) {
	var opt T
	valOpts := rf.ValueOf(opt)
	nh := func(w http.ResponseWriter, r *http.Request) {
		params := []any{}
		for i := 0; i < valOpts.NumField(); i++ {
			field := valOpts.Type().Field(i)
			name := valOpts.Type().Field(i).Name
			p := r.URL.Query().Get(name)

			var res any

			switch field.Type.Kind() {
			case rf.Int:
				res, _ = strconv.Atoi(p)
			default:
				res = p
			}

			params = append(params, res)

			for _, p := range params {
				switch p := p.(type) {
				case string:
					SetAttr(opt, name, p)
				case int:
					SetAttr(opt, name, p)
				}
			}

		}
		res := h(opt)
		w.Write(res.toJson())
	}

	Hdls = append(Hdls, Pat{Method: method, Path: path, Handler: nh})
}

type osch struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (o osch) toJson() []byte {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return data
}

func Tf() {
	type tOpt struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	thandle(http.MethodGet, "/", func(opt tOpt) OutSchema {
		return osch{
			Id:   opt.Id,
			Name: opt.Name,
		}
	})

	fmt.Println(Hdls)
}
