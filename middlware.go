package lax

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// Middleware is Go idiomatic middleware which logs request statistic data using lax logger.
// Additionally, If Chi's middleware.RequestID is used before this middleware, request id will be logged.
type Middleware struct {
	next http.Handler
	log  Logger
}

// NewMiddleware creates new lax logger middleware.
func NewMiddleware(next http.Handler, log Logger) *Middleware {
	return &Middleware{
		next: next,
		log:  log,
	}
}

// ServeHTTP implements http.Handler interface.
func (m *Middleware) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	writer := middleware.NewWrapResponseWriter(res, req.ProtoMajor)
	start := time.Now()

	m.next.ServeHTTP(writer, req)

	if reqID := middleware.GetReqID(req.Context()); reqID != "" {
		m.log.Info(
			"request completed",
			Time("time", start),
			String("req_id", reqID),
			String("method", req.Method),
			String("uri", req.RequestURI),
			Int("status", writer.Status()),
			Int("bytes", writer.BytesWritten()),
			Duration("duration", time.Since(start)))

		return
	}

	m.log.Info(
		"request completed",
		Time("time", start),
		String("method", req.Method),
		String("uri", req.RequestURI),
		Int("status", writer.Status()),
		Int("bytes", writer.BytesWritten()),
		Duration("duration", time.Since(start)))
}
