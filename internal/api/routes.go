package api

import (
	"net/http"
)

// Middleware is a function that takes  http.HandlerFunc and returns a new http.HandlerFunc.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// NewRouter returns a new router with an empty middleware stack.
func NewRouter() *Router {
	r := &Router{middlewares: []Middleware{}}
	r.mux = http.NewServeMux()
	return r
}

// Router is a simple router that supports middleware.
type Router struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// Use adds a new middleware to the middleware stack.
func (r *Router) Use(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}

// Get adds a new GET route to the router.
func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.mux.HandleFunc(path, r.applyMiddlewares(handler))
}

// Delete adds a new DELETE route to the router.
func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.mux.HandleFunc(path, r.applyMiddlewares(handler))
}

// Post adds a new POST route to the router.
func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.mux.HandleFunc(path, r.applyMiddlewares(handler))
}

// applyMiddlewares returns a new http.HandlerFunc that applies all the middlewares to the original handler.
func (r *Router) applyMiddlewares(handler http.HandlerFunc) http.HandlerFunc {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}
	return handler
}
