package app

import (
	"github.com/valyala/fasthttp"
)

const (
	contentTypeJSON = "application/json"
)

var (
	bodyNotFound = []byte(`{"code":404}`)
)

func (app *Application) HandleQuery(ctx *fasthttp.RequestCtx) {
	app.logger.Debugf("incoming request: %s %s", ctx.Method(), ctx.URI().String())
	resp, ok := app.index.QueryBytes(string(ctx.Path()))
	if !ok {
		ctx.SetContentType(contentTypeJSON)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBody(bodyNotFound)
		return
	}

	ctx.SetContentType(contentTypeJSON)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(resp)
}
