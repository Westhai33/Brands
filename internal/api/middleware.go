package api

import (
	"Brands/internal/config"
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	"Brands/pkg/zerohook"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"
)

func CORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Headers", config.CorsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", config.CorsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", config.CorsAllowOrigin)

		next(ctx)
	}
}

func RecoveryMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if rvr := recover(); rvr != nil {
				stack := debug.Stack()

				zerohook.Logger.Error().
					Interface("panic", rvr).
					Str("method", string(ctx.Method())).
					Str("url", string(ctx.URI().String())).
					Str("remote_addr", ctx.RemoteAddr().String()).
					Str("stack_trace", string(stack)).
					Msg("Recovered from panic")

				pc, file, line, ok := runtime.Caller(3)
				if ok {
					zerohook.Logger.Error().
						Str("file", file).
						Int("line", line).
						Str("function", runtime.FuncForPC(pc).Name()).
						Msg("Panic occurred here")
				}
				ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)

			}
		}()

		next(ctx)
	}
}

func TraceMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		requestID, ok := ctx.UserValue("request-id").(string)
		if !ok {
			requestID = uuid.New().String()
			ctx.SetUserValue("request-id", requestID)
		}
		ctxWithBackground, ok := ctx.UserValue("traceContext").(context.Context)
		if !ok {
			ctxWithBackground = context.Background()
		}
		operationName := fmt.Sprintf("%s %s", string(ctx.Method()), string(ctx.Path()))
		span, spanCtx := opentracing.StartSpanFromContext(ctxWithBackground, operationName)
		defer span.Finish()

		span.SetTag("request_id", requestID)
		span.SetTag("http.method", string(ctx.Method()))
		span.SetTag("http.url", string(ctx.URI().String()))
		span.SetTag("http.route", string(ctx.Path()))
		span.SetTag("host.name", string(ctx.Host()))
		span.SetTag("http.user_agent", string(ctx.Request.Header.UserAgent()))
		span.SetTag("http.scheme", string(ctx.Request.URI().Scheme()))
		span.SetTag("http.flavor", string(ctx.Request.Header.Protocol()))
		span.SetTag("host.ip", string(ctx.RemoteIP().String()))
		span.SetTag("http.request.body_size", fmt.Sprintf("%d", len(ctx.Request.Body())))
		span.SetTag("http.request.body_stringify", string(ctx.Request.Body()))

		zerohook.Logger.Debug().
			Str("request_id", requestID).
			Msg("Trace context set")

		ctx.SetUserValue("traceContext", spanCtx)
		next(ctx)

		span.SetTag("http.status_code", ctx.Response.StatusCode())
	}
}

// LoggingMiddleware логирует каждый запрос с включением request_id.
func LoggingMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		requestID, ok := ctx.UserValue("request-id").(string)
		if !ok {
			requestID = uuid.New().String()
			ctx.SetUserValue("request-id", requestID)
		}
		traceCtx, ok := ctx.UserValue("traceContext").(context.Context)
		if !ok {
			traceCtx = context.Background()
		}
		ctxWithRequestID := context.WithValue(traceCtx, "request-id", requestID)
		ctx.SetUserValue("traceContext", ctxWithRequestID)

		ctx.SetUserValue("traceContext", ctxWithRequestID)
		begin := time.Now()
		next(ctx)
		end := time.Now()
		zerohook.Logger.Info().
			Str("request_id", requestID).
			Bytes("method", ctx.Method()).
			Str("url", string(ctx.URI().String())).
			Int("status", ctx.Response.StatusCode()).
			Dur("latency", end.Sub(begin)).
			Msg("Completed request")

	}
}
