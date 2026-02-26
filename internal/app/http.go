package app

import "github.com/valyala/fasthttp"

const (
	contentTypeJSON = "application/json"
)

func (app *Application) HandleQuery(ctx *fasthttp.RequestCtx) {
	resp, ok := app.index.QueryBytes(string(ctx.Path()))
	if !ok {
		ctx.NotFound()
		return
	}

	ctx.SetContentType(contentTypeJSON)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(resp)
}
