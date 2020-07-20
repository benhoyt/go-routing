// Test various ways to do HTTP method+path routing in Go

/*

Implement the URLs below

1) switch or if-else based router with plain ServeMux
1a) switch with match() function?
2) custom tiny regexp router (like Gifty)
3) Axel Wagner's approach
4) one or two popular Go routers: chi and something else, maybe https://github.com/bmizerany/pat or gorilla/mux

GET  /										# home
GET  /contact								# contact

GET  /api/widgets							# apiGetWidgets
POST /api/widgets                           # apiCreateWidget
POST /api/widgets/:slug                     # apiUpdateWidget
POST /api/widgets/:slug/parts               # apiCreateWidgetPart
POST /api/widgets/:slug/parts/:id/update    # apiUpdateWidgetPart
POST /api/widgets/:slug/parts/:id/delete    # apiDeleteWidgetPart

GET  /:slug									# widget
GET  /:slug/admin                           # widgetAdmin
POST /:slug/image							# widgetImage

*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/benhoyt/go-routing/routematch"
	"github.com/benhoyt/go-routing/routesplit"
)

const port = 8080

func main() {
	if len(os.Args) < 2 || routers[os.Args[1]] == nil {
		routerNames := []string{}
		for k := range routers {
			routerNames = append(routerNames, k)
		}
		sort.Strings(routerNames)
		fmt.Fprintf(os.Stderr, "usage: go-routing router\n\n")
		fmt.Fprintf(os.Stderr, "router is one of: %s\n", strings.Join(routerNames, ", "))
		os.Exit(1)
	}
	routerName := os.Args[1]
	router := routers[routerName]

	http.Handle("/", router)
	fmt.Printf("listening on port %d using %s router\n", port, routerName)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

var routers = map[string]http.Handler{
	"match": http.HandlerFunc(routematch.Route),
	"split": http.HandlerFunc(routesplit.Route),
}

func routeMatch(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	var slug string
	var id int

	p := r.URL.Path
	switch {
	case match(p, "/"):
		h = get(home)
	case match(p, "/contact"):
		h = get(contact)
	case match(p, "/api/widgets") && isGet(r):
		h = get(apiGetWidgets)
	case match(p, "/api/widgets"):
		h = post(apiCreateWidget)
	case match(p, "/api/widgets/%s", &slug):
		h = post(apiWidget{slug}.update)
	case match(p, "/api/widgets/%s/parts", &slug):
		h = post(apiWidget{slug}.createPart)
	case match(p, "/api/widgets/%s/parts/%d/update", &slug, &id):
		h = post(apiWidgetPart{slug, id}.update)
	case match(p, "/api/widgets/%s/parts/%d/delete", &slug, &id):
		h = post(apiWidgetPart{slug, id}.delete)
	case match(p, "/%s", &slug):
		h = get(widget{slug}.widget)
	case match(p, "/%s/admin", &slug):
		h = get(widget{slug}.admin)
	case match(p, "/%s/image", &slug):
		h = post(widget{slug}.image)
	default:
		http.NotFound(w, r)
		return
	}
	h.ServeHTTP(w, r)
}

func match(path, pattern string, vars ...interface{}) bool {
	pathParts := strings.Split(path, "/")
	patternParts := strings.Split(pattern, "/")
	if len(pathParts) != len(patternParts) {
		return false
	}
	for i, pathPart := range pathParts {
		patternPart := patternParts[i]
		switch patternPart {
		case "%s":
			p := vars[0].(*string)
			vars = vars[1:]
			*p = pathPart
		case "%d":
			p := vars[0].(*int)
			vars = vars[1:]
			n, err := strconv.Atoi(pathPart)
			if err != nil || n <= 0 {
				return false
			}
			*p = n
		default:
			if pathPart != patternPart {
				return false
			}
		}
	}
	return true
}

func routeSplit(w http.ResponseWriter, r *http.Request) {
	// Split path into slash-separated parts, for example, path "/foo/bar"
	// gives p==["foo", "bar"] and path "/" gives p==[""].
	p := strings.Split(r.URL.Path, "/")[1:]
	n := len(p)

	var h http.Handler
	switch {
	case n == 1 && p[0] == "":
		h = get(home)
	case n == 1 && p[0] == "contact":
		h = get(contact)

	case n >= 2 && p[0] == "api" && p[1] == "widgets":
		// /api/widgets/*
		switch {
		case n == 2 && isGet(r):
			h = get(apiGetWidgets)
		case n == 2:
			h = post(apiCreateWidget)
		case n >= 3:
			// /api/widgets/:slug/*
			slug := p[2]
			var id int
			switch {
			case n == 3:
				h = post(apiWidget{slug}.update)
			case n == 4 && p[3] == "parts":
				h = post(apiWidget{slug}.createPart)
			case n == 6 && isId(p[4], &id):
				// /api/widgets/:slug/parts/:id/*
				switch {
				case p[5] == "update":
					h = post(apiWidgetPart{slug, id}.update)
				case p[5] == "delete":
					h = post(apiWidgetPart{slug, id}.delete)
				}
			}
		}

	case n >= 1:
		// /:slug/*
		slug := p[0]
		switch {
		case n == 1:
			h = get(widget{slug}.widget)
		case n == 2 && p[1] == "admin":
			h = get(widget{slug}.admin)
		case n == 2 && p[1] == "image":
			h = post(widget{slug}.image)
		}
	}

	if h == nil {
		http.NotFound(w, r)
		return
	}
	h.ServeHTTP(w, r)
}

func isGet(r *http.Request) bool {
	return r.Method == http.MethodGet || r.Method == http.MethodHead
}

func allowMethod(h http.HandlerFunc, methods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, m := range methods {
			if m == r.Method {
				h(w, r)
				return
			}
		}
		w.Header().Set("Allow", strings.Join(methods, ", "))
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func get(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, http.MethodGet, http.MethodHead)
}

func post(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, http.MethodPost)
}

func isId(s string, p *int) bool {
	id, err := strconv.Atoi(s)
	if err != nil || id <= 0 {
		return false
	}
	*p = id
	return true
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "home\n")
}

func contact(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "contact\n")
}

func apiGetWidgets(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "apiGetWidgets\n")
}

func apiCreateWidget(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "apiCreateWidget\n")
}

type apiWidget struct {
	slug string
}

func (h apiWidget) update(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiUpdateWidget %s\n", h.slug)
}

func (h apiWidget) createPart(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiCreateWidgetPart %s\n", h.slug)
}

type apiWidgetPart struct {
	slug string
	id   int
}

func (h apiWidgetPart) update(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiUpdateWidgetPart %s %d\n", h.slug, h.id)
}

func (h apiWidgetPart) delete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiDeleteWidgetPart %s %d\n", h.slug, h.id)
}

type widget struct {
	slug string
}

func (h widget) widget(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widget %s\n", h.slug)
}

func (h widget) admin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widgetAdmin %s\n", h.slug)
}

func (h widget) image(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widgetImage %s\n", h.slug)
}
