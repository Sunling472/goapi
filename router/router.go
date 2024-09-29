package router

import (
	"encoding/json"
	"log/slog"
	"net/http"
	rf "reflect"
	"strconv"
	"strings"

	"github.com/Sunling472/goapi/rftools"
	"github.com/Sunling472/goapi/types/methods"
)

type HandlerMeta struct {
	// TODO!
}

type Pattern struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
	Meta    map[string]any
}

type IRouter interface {
	GetMux() *http.ServeMux
	GetLog() *slog.Logger
	GetPatterns() []Pattern
	GetPath() string

	Include(rt IRouter)

	save(path, method string, handler http.HandlerFunc, meta map[string]any)
}

type Router struct {
	Path     string
	Mux      *http.ServeMux
	Log      *slog.Logger
	Patterns []Pattern
	Depends  []any
}

func New(path string) *Router {
	return &Router{
		Path:     path,
		Mux:      http.NewServeMux(),
		Log:      slog.Default(),
		Patterns: []Pattern{},
		Depends:  []any{},
	}
}

func (r Router) GetMux() *http.ServeMux {
	return r.Mux
}

func (r Router) GetLog() *slog.Logger {
	return r.Log
}

func (r Router) GetPatterns() []Pattern {
	return r.Patterns
}

func (r Router) GetPath() string {
	return r.Path
}

func (r *Router) save(path, method string, handler http.HandlerFunc, meta map[string]any) {
	r.Patterns = append(r.Patterns, Pattern{
		Method:  method,
		Path:    path,
		Handler: handler,
		Meta:    meta,
	})
}

func (r Router) Include(rt IRouter) {
	for _, p := range rt.GetPatterns() {
		r.Patterns = append(r.Patterns, p)
	}
}

type Schema interface {
	Json() []byte
}
type HandlerOpts[U, Q, In any] struct {
	Url      *U
	Query    *Q
	InSchema *In
}
type SmartHandlerFunc[U, Q, In any] func(opt HandlerOpts[U, Q, In]) Schema

func SmartHandler[M methods.METHOD, U, Q, In any](
	rt IRouter,
	path string,
	handler SmartHandlerFunc[U, Q, In],
) {
	const op = "goapi.router.SmartHandler"
	const (
		url   = "url"
		query = "query"
	)

	var (
		m          M
		urlParam   U
		queryParam Q
		inSchema   In
	)
	method := getMethod(m)

	resultMap := map[string]nameFieldMap{}
	urlVal := rf.ValueOf(urlParam)
	queryVal := rf.ValueOf(queryParam)

	urlMap, ok := getNameFieldMap(urlVal)
	if ok {
		resultMap[url] = urlMap
	}
	queryMap, ok := getNameFieldMap(queryVal)
	if ok {
		resultMap[query] = queryMap
	}

	setUrl := func(m nameFieldMap, r *http.Request) {
		for n, f := range m {
			p := r.PathValue(strings.ToLower(n))
			var res any

			switch f.Type.Kind() {
			case rf.Int:
				res, _ = strconv.Atoi(p)
				fallthrough
			default:
				if res == nil {
					res = p
				}
				switch p := res.(type) {
				default:
					rftools.SetAttr(&urlParam, n, p)
				}
			}
		}
	}

	setQuery := func(m nameFieldMap, r *http.Request) {
		for n, f := range m {
			v := r.URL.Query().Get(strings.ToLower(n))
			var res any

			switch f.Type.Kind() {
			case rf.Int:
				res, _ = strconv.Atoi(v)
				fallthrough
			default:
				res = v
				switch p := res.(type) {
				default:
					rftools.SetAttr(&queryParam, n, p)
				}
			}
		}
	}

	setBody := func(r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&inSchema)
		if err != nil {
			rt.GetLog().Info(
				"Validate Error",
				slog.String("op", op),
				slog.String("json-decode", err.Error()),
			)

		}
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		const (
			ctKey = "Content-Type"
			ctVal = "application/json"
		)
		w.Header().Set(ctKey, ctVal)

		for k, m := range resultMap {
			switch k {
			case url:
				setUrl(m, r)
			case query:
				setQuery(m, r)
			}
		}
		setBody(r)

		rt.GetLog().Info(
			"Request",
			slog.String("Method:", method),
			slog.String("Path:", path),
			slog.Any("UrlOpt:", urlParam),
			slog.Any("QueryOpt:", queryParam),
			slog.Any("Body:", inSchema),
		)

		opt := HandlerOpts[U, Q, In]{
			Url:      &urlParam,
			Query:    &queryParam,
			InSchema: &inSchema,
		}
		res := handler(opt)

		rt.GetLog().Info("Response", slog.Any("Json", res))

		w.Write(res.Json())
	}

	if rt.GetPath() != "" {
		path = rt.GetPath() + path
	}

	rt.save(path, method, h, map[string]any{})
}

func getMethod(m any) string {
	mt := rf.TypeOf(m)
	return strings.ToUpper(mt.Name())
}

type nameFieldMap map[string]rf.StructField

func getNameFieldMap(v rf.Value) (nameFieldMap, bool) {
	r := nameFieldMap{}
	if !v.IsValid() {
		return r, false
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		n := f.Name
		r[n] = f
	}
	if len(r) == 0 {
		return r, false
	}

	return r, true
}
