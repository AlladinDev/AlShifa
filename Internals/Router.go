package internals

import (
	"net/http"
)

// Middleware type
type Middleware func(http.Handler) http.Handler

// Router wraps http.ServeMux
type Router struct {
	mux *http.ServeMux
}

// NewRouter creates a new Router
func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

// ServeHTTP makes Router implement http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// generic method registration
func (r *Router) Handle(
	method string,
	path string,
	handler http.HandlerFunc,
	middlewares ...Middleware,
) {

	// Start with final handler
	var h http.Handler = handler

	// Wrap middlewares in reverse order for correct execution
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	// Wrap method guard
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		h.ServeHTTP(w, req)
	})

	// Register the final wrapped handler
	r.mux.Handle(path, finalHandler)
}

// Convenience methods for each HTTP verb
func (r *Router) GET(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodGet, path, handler, middlewares...)
}

func (r *Router) POST(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodPost, path, handler, middlewares...)
}

func (r *Router) PUT(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodPut, path, handler, middlewares...)
}

func (r *Router) DELETE(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodDelete, path, handler, middlewares...)
}

func (r *Router) PATCH(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodPatch, path, handler, middlewares...)
}

func (r *Router) OPTIONS(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodOptions, path, handler, middlewares...)
}

func (r *Router) HEAD(path string, handler http.HandlerFunc, middlewares ...Middleware) {
	r.Handle(http.MethodHead, path, handler, middlewares...)
}
