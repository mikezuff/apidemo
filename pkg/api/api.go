package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mikezuff/apidemo/pkg/config"
	"github.com/ziflex/lecho/v3"
)

func Run() {
	e := echo.New()
	e.HideBanner = true

	appContext := config.InitAppContext()

	e.Use(lecho.Middleware(
		lecho.Config{
			Logger: lecho.New(appContext.Logger),
			/*
				Enricher: func(c echo.Context, ctx zerolog.Context) zerolog.Context {
					if c.Request().Body != nil {
						requestBytes, err := ioutil.ReadAll(c.Request().Body)
						if err != nil {
							// this error will just disappear into the logs
							// so we're expecting a malformed body to cause an application error later?
							ctx.Err(fmt.Errorf("reading original request body: %w", err))
						} else {
							// reset the request body
							c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(requestBytes))
							ctx.Str("request_body", string(requestBytes))
						}
					}
					return ctx
				},
			*/
			HandleError: true,
		}))

	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		appContext.Logger.Info().
			Str("request_body", string(reqBody)).
			Str("response_body", string(resBody)).
			Msg("request/response body")
	}))

	recoverConfig := middleware.DefaultRecoverConfig
	recoverConfig.DisableStackAll = true
	e.Use(middleware.RecoverWithConfig(recoverConfig))
	e.Use(middleware.CORS())

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 10 * time.Second,
	}))

	// TODO: add prometheus
	// TODO: add rate limiting

	RegisterRoutes(e)

	// Start server
	// Use graceful shutdown: https://echo.labstack.com/cookbook/graceful-shutdown/
	go func() {
		address := ":8080"
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server; ", err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func RegisterRoutes(e *echo.Echo) {
	e.GET("/panic", panicHandler)
	e.GET("/clientip", clientIPHandler)
	e.GET("/remoteaddr", remoteAddrHandler)
}

func panicHandler(c echo.Context) error {
	panic("panic test")
}

func clientIPHandler(c echo.Context) error {
	return c.String(http.StatusOK, c.RealIP())
}

func remoteAddrHandler(c echo.Context) error {
	return c.String(http.StatusOK, c.Request().RemoteAddr)
}
