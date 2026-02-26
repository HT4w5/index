package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/HT4w5/index/internal/config"
	"github.com/HT4w5/index/pkg/index"
	"github.com/valyala/fasthttp"
)

type Application struct {
	cfg config.Config

	index   *index.Index
	httpsrv *fasthttp.Server
}

func New(cfg config.Config) *Application {
	return &Application{
		cfg: cfg,
	}
}

func (app *Application) Start() error {
	// Create index
	var err error
	opts := make([]func(*index.Index), 0)
	if app.cfg.Filesystem.Root != "" {
		opts = append(opts, index.WithRoot(app.cfg.Filesystem.Root))
	}
	if app.cfg.Cache.TTL != 0 {
		opts = append(opts, index.WithTTL(time.Duration(app.cfg.Cache.TTL)))
	}
	if app.cfg.Cache.MaxSize != 0 {
		opts = append(opts, index.WithMaxSize(app.cfg.Cache.MaxSize))
	}
	var level index.LogLevel
	switch app.cfg.Log.Level {
	case "none":
		level = index.None
	case "error":
		level = index.Error
	case "warn":
		level = index.Warn
	case "":
		fallthrough
	case "info":
		level = index.Info
	case "debug":
		level = index.Debug
	}
	opts = append(opts, index.WithLogger(&index.SimpleLogger{
		Level: level,
	}))

	app.index, err = index.New(opts...)
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}

	// HTTP listen
	app.httpsrv = &fasthttp.Server{
		Handler: app.HandleQuery,
	}

	addr := app.cfg.HTTP.Addr
	port := app.cfg.HTTP.Port
	if len(addr) == 0 {
		addr = "[::]"
	}
	if port == 0 {
		port = 80
	}

	go app.httpsrv.ListenAndServe(fmt.Sprintf("%s:%d", addr, port))

	return nil
}

func (app *Application) Shutdown() error {
	// HTTP shutdown
	err := app.httpsrv.Shutdown()

	// Index close
	return errors.Join(err, app.index.Close())
}
