// Go HTTP router based on strings.Split() with a switch statement

package split

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func Serve(w http.ResponseWriter, r *http.Request) {
	// Split path into slash-separated parts, for example, path "/foo/bar"
	// gives p==["foo", "bar"] and path "/" gives p==[""].
	p := strings.Split(r.URL.Path, "/")[1:]
	n := len(p)

	var h http.Handler
	var id int
	switch {
	case n == 1 && p[0] == "":
		h = get(home)
	case n == 1 && p[0] == "contact":
		h = get(contact)
	case n == 2 && p[0] == "api" && p[1] == "widgets" && r.Method == "GET":
		h = get(apiGetWidgets)
	case n == 2 && p[0] == "api" && p[1] == "widgets":
		h = post(apiCreateWidget)
	case n == 3 && p[0] == "api" && p[1] == "widgets" && p[2] != "":
		h = post(apiWidget{p[2]}.update)
	case n == 4 && p[0] == "api" && p[1] == "widgets" && p[2] != "" && p[3] == "parts":
		h = post(apiWidget{p[2]}.createPart)
	case n == 6 && p[0] == "api" && p[1] == "widgets" && p[2] != "" && p[3] == "parts" && isId(p[4], &id) && p[5] == "update":
		h = post(apiWidgetPart{p[2], id}.update)
	case n == 6 && p[0] == "api" && p[1] == "widgets" && p[2] != "" && p[3] == "parts" && isId(p[4], &id) && p[5] == "delete":
		h = post(apiWidgetPart{p[2], id}.delete)
	case n == 1:
		h = get(widget{p[0]}.widget)
	case n == 2 && p[1] == "admin":
		h = get(widget{p[0]}.admin)
	case n == 2 && p[1] == "image":
		h = post(widget{p[0]}.image)
	default:
		http.NotFound(w, r)
		return
	}
	h.ServeHTTP(w, r)
}

func allowMethod(h http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			w.Header().Set("Allow", method)
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}

func get(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, "GET")
}

func post(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, "POST")
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
