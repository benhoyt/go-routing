// Go HTTP router based on a simple custom match() function

package match

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func Serve(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	var apiWidget apiWidget
	var apiWidgetPart apiWidgetPart
	var widget widget

	p := r.URL.Path
	switch {
	case match(p, "/"):
		h = get(home)
	case match(p, "/contact"):
		h = get(contact)
	case match(p, "/api/widgets"):
		if r.Method == "GET" {
			h = get(apiGetWidgets)
		} else {
			h = post(apiCreateWidget)
		}
	case match(p, "/api/widgets/+", &apiWidget.slug):
		h = post(withURLParams(apiWidgetUpdate, apiWidget))
	case match(p, "/api/widgets/+/parts", &apiWidget.slug):
		h = post(withURLParams(apiWidgetCreatePart, apiWidget))
	case match(p, "/api/widgets/+/parts/+/update", &apiWidgetPart.slug, &apiWidgetPart.id):
		h = post(withURLParams(apiWidgetPartUpdate, apiWidgetPart))
	case match(p, "/api/widgets/+/parts/+/delete", &apiWidgetPart.slug, &apiWidgetPart.id):
		h = post(withURLParams(apiWidgetPartDelete, apiWidgetPart))
	case match(p, "/+", &widget.slug):
		h = get(withURLParams(widgetWidget, widget))
	case match(p, "/+/admin", &widget.slug):
		h = get(withURLParams(widgetAdmin, widget))
	case match(p, "/+/image", &widget.slug):
		h = post(withURLParams(widgetImage, widget))
	default:
		http.NotFound(w, r)
		return
	}
	h.ServeHTTP(w, r)
}

// match reports whether path matches the given pattern, which is a
// path with '+' wildcards wherever you want to use a parameter. Path
// parameters are assigned to the pointers in vars (len(vars) must be
// the number of wildcards), which must be of type *string or *int.
func match(path, pattern string, vars ...interface{}) bool {
	for ; pattern != "" && path != ""; pattern = pattern[1:] {
		switch pattern[0] {
		case '+':
			// '+' matches till next slash in path
			slash := strings.IndexByte(path, '/')
			if slash < 0 {
				slash = len(path)
			}
			segment := path[:slash]
			path = path[slash:]
			switch p := vars[0].(type) {
			case *string:
				*p = segment
			case *int:
				n, err := strconv.Atoi(segment)
				if err != nil || n < 0 {
					return false
				}
				*p = n
			default:
				panic("vars must be *string or *int")
			}
			vars = vars[1:]
		case path[0]:
			// non-'+' pattern byte must match path byte
			path = path[1:]
		default:
			return false
		}
	}
	return path == "" && pattern == ""
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

type ctxKey struct{}

func withURLParams(h http.HandlerFunc, params interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxKey{}, params)
		h(w, r.WithContext(ctx))
	}
}

func urlParams(r *http.Request) interface{} {
	return r.Context().Value(ctxKey{})
}

func get(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, "GET")
}

func post(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, "POST")
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

func apiWidgetUpdate(w http.ResponseWriter, r *http.Request) {
	params := urlParams(r).(apiWidget)
	fmt.Fprintf(w, "apiUpdateWidget %s\n", params.slug)
}

func apiWidgetCreatePart(w http.ResponseWriter, r *http.Request) {
	params := urlParams(r).(apiWidget)
	fmt.Fprintf(w, "apiCreateWidgetPart %s\n", params.slug)
}

type apiWidgetPart struct {
	slug string
	id   int
}

func apiWidgetPartUpdate(w http.ResponseWriter, r *http.Request) {
	params := urlParams(r).(apiWidgetPart)
	fmt.Fprintf(w, "apiUpdateWidgetPart %s %d\n", params.slug, params.id)
}

func apiWidgetPartDelete(w http.ResponseWriter, r *http.Request) {
	params := urlParams(r).(apiWidgetPart)
	fmt.Fprintf(w, "apiDeleteWidgetPart %s %d\n", params.slug, params.id)
}

type widget struct {
	slug string
}

func widgetWidget(w http.ResponseWriter, r *http.Request) {
	params := urlParams(r).(widget)
	fmt.Fprintf(w, "widget %s\n", params.slug)
}

func widgetAdmin(w http.ResponseWriter, r *http.Request) {
	params := urlParams(r).(widget)
	fmt.Fprintf(w, "widgetAdmin %s\n", params.slug)
}

func widgetImage(w http.ResponseWriter, r *http.Request) {
	params := urlParams(r).(widget)
	fmt.Fprintf(w, "widgetImage %s\n", params.slug)
}
