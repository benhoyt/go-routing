// Routing based on Go 1.22's http.ServeMux enhancements

//go:build go1.22

package stdlib

import (
	"fmt"
	"net/http"
	"strconv"
)

var Serve http.Handler

func init() {
	r := http.NewServeMux()

	r.HandleFunc("GET /{$}", home)
	r.HandleFunc("GET /contact", contact)
	r.HandleFunc("GET /api/widgets", apiGetWidgets)
	r.HandleFunc("POST /api/widgets", apiCreateWidget)
	r.HandleFunc("POST /api/widgets/{slug}", apiUpdateWidget)
	r.HandleFunc("POST /api/widgets/{slug}/parts", apiCreateWidgetPart)
	r.HandleFunc("POST /api/widgets/{slug}/parts/{id}/update", apiUpdateWidgetPart)
	r.HandleFunc("POST /api/widgets/{slug}/parts/{id}/delete", apiDeleteWidgetPart)
	r.HandleFunc("GET /{slug}", widgetGet)
	r.HandleFunc("GET /{slug}/admin", widgetAdmin)
	r.HandleFunc("POST /{slug}/image", widgetImage)

	Serve = r
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

func apiUpdateWidget(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	fmt.Fprintf(w, "apiUpdateWidget %s\n", slug)
}

func apiCreateWidgetPart(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	fmt.Fprintf(w, "apiCreateWidgetPart %s\n", slug)
}

func apiUpdateWidgetPart(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "apiUpdateWidgetPart %s %d\n", slug, id)
}

func apiDeleteWidgetPart(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "apiDeleteWidgetPart %s %d\n", slug, id)
}

func widgetGet(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	fmt.Fprintf(w, "widget %s\n", slug)
}

func widgetAdmin(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	fmt.Fprintf(w, "widgetAdmin %s\n", slug)
}

func widgetImage(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	fmt.Fprintf(w, "widgetImage %s\n", slug)
}
