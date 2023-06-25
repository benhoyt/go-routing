// Routing based on github.com/jba/muxpatterns (proposed enhancements to http.ServeMux)

package muxpatterns

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jba/muxpatterns"
)

var (
	Serve http.Handler

	// Wouldn't use a global for this in "real life" -- this is only needed
	// because we have to call mux.PathValue(r, "name") below instead of
	// r.PathValue("name") in the actual proposal.
	mux *muxpatterns.ServeMux
)

func init() {
	r := muxpatterns.NewServeMux()

	r.HandleFunc("GET /{$}", home)
	r.HandleFunc("GET /contact", contact)
	r.HandleFunc("GET /api/widgets", apiGetWidgets)
	r.HandleFunc("POST /api/widgets", apiCreateWidget)
	// TODO: seems to have a bug where "POST /api/widgets/" matches and gives slug "/"
	r.HandleFunc("POST /api/widgets/{slug}", apiUpdateWidget)
	r.HandleFunc("POST /api/widgets/{slug}/parts", apiCreateWidgetPart)
	r.HandleFunc("POST /api/widgets/{slug}/parts/{id}/update", apiUpdateWidgetPart)
	r.HandleFunc("POST /api/widgets/{slug}/parts/{id}/delete", apiDeleteWidgetPart)
	r.HandleFunc("GET /{slug}", widgetGet)
	r.HandleFunc("GET /{slug}/admin", widgetAdmin)
	r.HandleFunc("POST /{slug}/image", widgetImage)

	mux = r
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
	slug := mux.PathValue(r, "slug")
	fmt.Fprintf(w, "apiUpdateWidget %s\n", slug)
}

func apiCreateWidgetPart(w http.ResponseWriter, r *http.Request) {
	slug := mux.PathValue(r, "slug")
	fmt.Fprintf(w, "apiCreateWidgetPart %s\n", slug)
}

func apiUpdateWidgetPart(w http.ResponseWriter, r *http.Request) {
	slug := mux.PathValue(r, "slug")
	idStr := mux.PathValue(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "apiUpdateWidgetPart %s %d\n", slug, id)
}

func apiDeleteWidgetPart(w http.ResponseWriter, r *http.Request) {
	slug := mux.PathValue(r, "slug")
	idStr := mux.PathValue(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "apiDeleteWidgetPart %s %d\n", slug, id)
}

func widgetGet(w http.ResponseWriter, r *http.Request) {
	slug := mux.PathValue(r, "slug")
	fmt.Fprintf(w, "widget %s\n", slug)
}

func widgetAdmin(w http.ResponseWriter, r *http.Request) {
	slug := mux.PathValue(r, "slug")
	fmt.Fprintf(w, "widgetAdmin %s\n", slug)
}

func widgetImage(w http.ResponseWriter, r *http.Request) {
	slug := mux.PathValue(r, "slug")
	fmt.Fprintf(w, "widgetImage %s\n", slug)
}
