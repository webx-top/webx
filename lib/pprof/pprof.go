package pprof

import (
	"net/http/pprof"

	"github.com/labstack/echo"
)

// Wrap adds several routes from package `net/http/pprof` to *gin.Engine object
func Wrap(router *echo.Echo) {
	router.Get("/debug/pprof/", IndexHandler())
	router.Get("/debug/pprof/heap", HeapHandler())
	router.Get("/debug/pprof/goroutine", GoroutineHandler())
	router.Get("/debug/pprof/block", BlockHandler())
	router.Get("/debug/pprof/threadcreate", ThreadCreateHandler())
	router.Get("/debug/pprof/cmdline", CmdlineHandler())
	router.Get("/debug/pprof/profile", ProfileHandler())
	router.Get("/debug/pprof/symbol", SymbolHandler())
	router.Get("/debug/pprof/trace", TraceHandler())
}

// Wrapper make sure we are backward compatible
var Wrapper = Wrap

// IndexHandler will pass the call from /debug/pprof to pprof
func IndexHandler() echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		pprof.Index(ctx.Response(), ctx.Request())
		return nil
	}
}

// HeapHandler will pass the call from /debug/pprof/heap to pprof
func HeapHandler() echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		pprof.Handler("heap").ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	}
}

// GoroutineHandler will pass the call from /debug/pprof/goroutine to pprof
func GoroutineHandler() echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		pprof.Handler("goroutine").ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	}
}

// BlockHandler will pass the call from /debug/pprof/block to pprof
func BlockHandler() echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		pprof.Handler("block").ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	}
}

// ThreadCreateHandler will pass the call from /debug/pprof/threadcreate to pprof
func ThreadCreateHandler() echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		pprof.Handler("threadcreate").ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	}
}

// CmdlineHandler will pass the call from /debug/pprof/cmdline to pprof
func CmdlineHandler() echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		pprof.Cmdline(ctx.Response(), ctx.Request())
		return nil
	}
}

// ProfileHandler will pass the call from /debug/pprof/profile to pprof
func ProfileHandler() echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		pprof.Profile(ctx.Response(), ctx.Request())
		return nil
	}
}

// SymbolHandler will pass the call from /debug/pprof/symbol to pprof
func SymbolHandler() echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		pprof.Symbol(ctx.Response(), ctx.Request())
		return nil
	}
}

// TraceHandler will pass the call from /debug/pprof/trace to pprof
func TraceHandler() echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		pprof.Trace(ctx.Response(), ctx.Request())
		return nil
	}
}
