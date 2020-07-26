// Routing based on Axel Wagner's ShiftPath approach:
// https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html

package shiftpath

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
)

func Serve(w http.ResponseWriter, r *http.Request) {
	// Pre-emptively return Not Found for URLs with trailing slash,
	// as ShiftPath approach doesn't distinguish between no trailing
	// slash and trailing slash
	if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	}

	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "":
		get(home)(w, r)
	case "api":
		api(w, r)
	case "contact":
		get(contact)(w, r)
	default:
		widget{head}.root(w, r)
	}
}

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
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
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head != "" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "contact\n")
}

func api(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "widgets":
		apiWidgets(w, r)
	default:
		http.NotFound(w, r)
	}
}

func apiWidgets(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "":
		if r.Method == "GET" {
			apiGetWidgets(w, r)
		} else {
			post(apiCreateWidget)(w, r)
		}
	default:
		apiWidget{head}.root(w, r)
	}
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

func (h apiWidget) root(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "":
		post(h.update)(w, r)
	case "parts":
		h.parts(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h apiWidget) update(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiUpdateWidget %s\n", h.slug)
}

func (h apiWidget) parts(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "":
		post(h.createPart)(w, r)
	default:
		id, err := strconv.Atoi(head)
		if err != nil || id <= 0 {
			http.NotFound(w, r)
			return
		}
		apiWidgetPart{h.slug, id}.root(w, r)
	}
}

func (h apiWidget) createPart(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiCreateWidgetPart %s\n", h.slug)
}

type apiWidgetPart struct {
	slug string
	id   int
}

func (h apiWidgetPart) root(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "update":
		post(h.update)(w, r)
	case "delete":
		post(h.delete)(w, r)
	default:
		http.NotFound(w, r)
	}
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

func (h widget) root(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "":
		get(h.get)(w, r)
	case "admin":
		get(h.admin)(w, r)
	case "image":
		post(h.image)(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h widget) get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widget %s\n", h.slug)
}

func (h widget) admin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widgetAdmin %s\n", h.slug)
}

func (h widget) image(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widgetImage %s\n", h.slug)
}
