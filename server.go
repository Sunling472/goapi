package goapi

import (
	"log/slog"
	"net/http"
	"os"
	"os/exec"

	"github.com/Sunling472/goapi/pkg/oapi"
	"github.com/Sunling472/goapi/router"
)

type ApiConfig struct {
	Info oapi.Info
	Path oapi.Path
}

type GoAPI struct {
	Router    router.IRouter
	ApiConfig *ApiConfig
	Server    *http.Server
}

func New(rt router.IRouter, cfg *ApiConfig) *GoAPI {
	return &GoAPI{
		Router:    rt,
		ApiConfig: cfg,
		Server: &http.Server{
			Addr:    "localhost:8080",
			Handler: rt.GetMux(),
		},
	}
}

func (g *GoAPI) registerHandlers() {
	for _, p := range g.Router.GetPatterns() {
		g.Router.GetMux().Handle(p.Method+" "+p.Path, p.Handler)
		g.Router.GetLog().Info(
			"Register",
			slog.String("Method", p.Method),
			slog.String("path", p.Path),
		)
	}
}

func Generate(c *ApiConfig) {
	const filename = "openapi.json"
	s := oapi.OpenAPI{
		OpenApi: "0.3.1",
		Info:    c.Info,
		Paths:   c.Path,
	}
	data, err := s.ToJson()
	if err != nil {
		panic(err)
	}

	nf := exec.Command("touch", "./"+filename)
	err = nf.Start()
	if err != nil {
		panic(err)
	}

	fp, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fp = fp + "/" + filename

	f, err := os.OpenFile(fp, 0666, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.Write(data)
}

func (g GoAPI) Serve() {
	g.registerHandlers()
	g.Router.GetLog().Info(
		"GoAPI Serve",
		"url", "http://"+g.Server.Addr,
	)

	if err := g.Server.ListenAndServe(); err != nil {
		g.Router.GetLog().Error("Serve Error", "msg", err.Error())
	}
}
