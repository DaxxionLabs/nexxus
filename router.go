package nexxus

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	INT    PathVarType = "int"
	STRING PathVarType = "string"
	UUID   PathVarType = "uuid"
)

type Router struct {
	// List of routes in the system
	routes []*Route

	// all the middleware that has to be called in a router
	middlewares []middleware

	// names routes to be used for quick lookup for template stuff later on
	namedRoutes map[string]*Route

	// default handler for 404 returns
	notFoundHandler Handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	newCtx, cancel := context.WithTimeout(req.Context(), 60*time.Second)

	ctx := Context{
		writer:  w,
		request: req,
		context: newCtx,
		cancel:  cancel,
	}

	var handler Handler
	var matched MatchedRoute
	if r.canServe(req, &matched) {
		handler = matched.Handler
		if err := handler(ctx); err != nil {
			// TODO: Handle the errors here
		}
	} else {
		_ = r.notFoundHandler(ctx)
	}
}

func (r *Router) Use(middleware ...middleware) {
	r.middlewares = append(r.middlewares, middleware...)
}

func (r *Router) canServe(req *http.Request, matched *MatchedRoute) bool {
	for _, rt := range r.routes {
		if rt.Path == req.URL.Path && rt.Method == req.Method {

			// setup router middleware chain
			for i := len(r.middlewares) - 1; i >= 0; i-- {
				matched.Handler = r.middlewares[i].Middleware(rt.Handler)
			}

			// setup route level middleware chain
			for i := len(rt.middlewares) - 1; i >= 0; i-- {
				matched.Handler = rt.middlewares[i].Middleware(matched.Handler)
			}

			matched.Path = rt.Path
			matched.Method = rt.Method
			return true
		}
	}
	return false
}

func NewRouter() *Router {
	return &Router{namedRoutes: make(map[string]*Route)}
}

func namedRoute(name, method, path string, handler Handler, middleware ...middleware) *Route {
	pathVars, err := parsePathVars(path)
	if err != nil {
		log.Fatal(err)
	}
	r := &Route{
		Name:        name,
		Method:      method,
		Path:        path,
		Handler:     handler,
		Vars:        pathVars,
		middlewares: middleware,
	}
	return r
}

func parsePathVars(path string) ([]PathVar, error) {
	vars := make([]PathVar, 0)
	pathBytes := []byte(path)

	for i := 0; i < len(pathBytes); i++ {
		if pathBytes[i] == '{' {
			rv, count, err := parseKeyAndType(pathBytes[i+1:])
			if err != nil {
				return nil, err
			}

			vars = append(vars, rv)
			i += count
		}
	}
	return vars, nil
}

func parseKeyAndType(start []byte) (PathVar, int, error) {
	end := false
	pathVar := make([]byte, 0)

	count := 0
	for i := 0; i < len(start); i++ {
		if start[i] == '}' {
			end = true
			count++
			break
		}
		pathVar = append(pathVar, start[i])
		count++
	}

	if !end {
		return PathVar{}, 0, errors.New(fmt.Sprintf("invalid path variable format :: %s", string(start)))
	}

	typeName := strings.Split(string(pathVar), ":")
	if len(typeName) != 2 {
		return PathVar{}, 0, errors.New(fmt.Sprintf("invalid path variable format found in path :: %s", string(pathVar[1:])))
	}

	varType := PathVarType(typeName[1])
	if varType != INT && varType != STRING && varType != UUID {
		return PathVar{}, 0, errors.New(fmt.Sprintf("invalid type \"%s\" in path variable declaration :: %s", typeName[1], string(pathVar[1:])))
	}

	return PathVar{name: typeName[0], valueType: varType}, count, nil
}

type MatchedRoute struct {
	Path    string
	Method  string
	Handler Handler
}

type PathVarType string
type PathVar struct {
	name      string
	valueType PathVarType
}

type Route struct {
	Path        string
	Method      string
	Handler     Handler
	Name        string
	Vars        []PathVar
	middlewares []middleware
}

func (r *Route) With(middleware ...middleware) {
	r.middlewares = append(r.middlewares, middleware...)
}
