// Go HTTP router based on a simple custom match() function

package routematch

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func Route(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	var slug string
	var id int

	p := r.URL.Path
	switch {
	case match(p, "/"):
		h = get(home)
	case match(p, "/contact"):
		h = get(contact)
	case match(p, "/api/widgets") && r.Method == "GET":
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
			if pathPart == "" {
				return false
			}
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

func allowMethod(h http.HandlerFunc, methods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, m := range methods {
			if m == r.Method {
				h(w, r)
				return
			}
		}
		w.Header().Set("Allow", strings.Join(methods, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	}
}

func get(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, http.MethodGet, http.MethodHead)
}

func post(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, http.MethodPost)
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
