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
		serveHome(w, r)
	case "api":
		serveApi(w, r)
	case "contact":
		serveContact(w, r)
	default:
		widget{head}.ServeHTTP(w, r)
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

func allowMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if method != r.Method {
		w.Header().Set("Allow", method)
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, "GET") {
		return
	}
	fmt.Fprint(w, "home\n")
}

func serveContact(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head != "" {
		http.NotFound(w, r)
		return
	}
	if !allowMethod(w, r, "GET") {
		return
	}
	fmt.Fprint(w, "contact\n")
}

func serveApi(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "widgets":
		serveApiWidgets(w, r)
	default:
		http.NotFound(w, r)
	}
}

func serveApiWidgets(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "":
		if r.Method == "GET" {
			serveApiGetWidgets(w, r)
		} else {
			serveApiCreateWidget(w, r)
		}
	default:
		apiWidget{head}.ServeHTTP(w, r)
	}
}

func serveApiGetWidgets(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, "GET") {
		return
	}
	fmt.Fprint(w, "apiGetWidgets\n")
}

func serveApiCreateWidget(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, "POST") {
		return
	}
	fmt.Fprint(w, "apiCreateWidget\n")
}

type apiWidget struct {
	slug string
}

func (h apiWidget) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "":
		h.serveUpdate(w, r)
	case "parts":
		h.serveParts(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h apiWidget) serveUpdate(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, "POST") {
		return
	}
	fmt.Fprintf(w, "apiUpdateWidget %s\n", h.slug)
}

func (h apiWidget) serveParts(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "":
		h.serveCreatePart(w, r)
	default:
		id, err := strconv.Atoi(head)
		if err != nil || id <= 0 {
			http.NotFound(w, r)
			return
		}
		apiWidgetPart{h.slug, id}.ServeHTTP(w, r)
	}
}

func (h apiWidget) serveCreatePart(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, "POST") {
		return
	}
	fmt.Fprintf(w, "apiCreateWidgetPart %s\n", h.slug)
}

type apiWidgetPart struct {
	slug string
	id   int
}

func (h apiWidgetPart) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "update":
		h.serveUpdate(w, r)
	case "delete":
		h.serveDelete(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h apiWidgetPart) serveUpdate(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, "POST") {
		return
	}
	fmt.Fprintf(w, "apiUpdateWidgetPart %s %d\n", h.slug, h.id)
}

func (h apiWidgetPart) serveDelete(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, "POST") {
		return
	}
	fmt.Fprintf(w, "apiDeleteWidgetPart %s %d\n", h.slug, h.id)
}

type widget struct {
	slug string
}

func (h widget) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "":
		h.serveGet(w, r)
	case "admin":
		h.serveAdmin(w, r)
	case "image":
		h.serveUpdateImage(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h widget) serveGet(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, "GET") {
		return
	}
	fmt.Fprintf(w, "widget %s\n", h.slug)
}

func (h widget) serveAdmin(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, "GET") {
		return
	}
	fmt.Fprintf(w, "widgetAdmin %s\n", h.slug)
}

func (h widget) serveUpdateImage(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, "POST") {
		return
	}
	fmt.Fprintf(w, "widgetImage %s\n", h.slug)
}
