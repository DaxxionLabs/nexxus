package nexxus

import (
	"context"
	"net/http"
)

type Context struct {
	Params   map[string]string
	PathVars map[string]string

	writer  http.ResponseWriter
	request *http.Request
	context context.Context
	cancel  context.CancelFunc
}

type Handler func(ctx Context) error

type HandlerFunc func(path string, handler func(w http.ResponseWriter, r *http.Request))

type middleware interface {
	Middleware(Handler) Handler
}

type MiddlewareFunc func(Handler) Handler

func (mw MiddlewareFunc) Middleware(h Handler) Handler {
	return mw(h)
}

type Nexxus struct {
	router *Router
	server *http.Server
}

func New() *Nexxus {
	router := NewRouter()
	return &Nexxus{
		router: router,
		server: &http.Server{Handler: router},
	}
}

func (n *Nexxus) Start() {}

func (n *Nexxus) Delete(name, path string, handler Handler) *Route {
	r := namedRoute(name, http.MethodDelete, path, handler)
	n.router.routes = append(n.router.routes, r)
	return r
}

func (n *Nexxus) Get(name, path string, handler Handler) *Route {
	r := namedRoute(name, http.MethodGet, path, handler)
	n.router.routes = append(n.router.routes, r)
	return r
}

func (n *Nexxus) Post(name, path string, handler Handler) *Route {
	r := namedRoute(name, http.MethodPost, path, handler)
	n.router.routes = append(n.router.routes, r)
	return r
}

func (n *Nexxus) Patch(name, path string, handler Handler) *Route {
	r := namedRoute(name, http.MethodPatch, path, handler)
	n.router.routes = append(n.router.routes, r)
	return r
}

func (n *Nexxus) Put(name, path string, handler Handler) *Route {
	r := namedRoute(name, http.MethodPost, path, handler)
	n.router.routes = append(n.router.routes, r)
	return r
}
