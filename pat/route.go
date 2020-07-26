// Routing based on the bmizerany/pat router

package pat

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bmizerany/pat"
)

var Serve http.Handler

func init() {
	r := pat.New()

	r.Get("/", http.HandlerFunc(home))
	r.Get("/contact", http.HandlerFunc(contact))
	r.Get("/api/widgets", http.HandlerFunc(apiGetWidgets))
	r.Post("/api/widgets", http.HandlerFunc(apiCreateWidget))
	r.Post("/api/widgets/:slug", http.HandlerFunc(apiUpdateWidget))
	r.Post("/api/widgets/:slug/parts", http.HandlerFunc(apiCreateWidgetPart))
	r.Post("/api/widgets/:slug/parts/:id/update", http.HandlerFunc(apiUpdateWidgetPart))
	r.Post("/api/widgets/:slug/parts/:id/delete", http.HandlerFunc(apiDeleteWidgetPart))
	r.Get("/:slug", http.HandlerFunc(widgetGet))
	r.Get("/:slug/admin", http.HandlerFunc(widgetAdmin))
	r.Post("/:slug/image", http.HandlerFunc(widgetImage))

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
	slug := r.URL.Query().Get(":slug")
	fmt.Fprintf(w, "apiUpdateWidget %s\n", slug)
}

func apiCreateWidgetPart(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get(":slug")
	fmt.Fprintf(w, "apiCreateWidgetPart %s\n", slug)
}

func apiUpdateWidgetPart(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get(":slug")
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "apiUpdateWidgetPart %s %d\n", slug, id)
}

func apiDeleteWidgetPart(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get(":slug")
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "apiDeleteWidgetPart %s %d\n", slug, id)
}

func widgetGet(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get(":slug")
	fmt.Fprintf(w, "widget %s\n", slug)
}

func widgetAdmin(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get(":slug")
	fmt.Fprintf(w, "widgetAdmin %s\n", slug)
}

func widgetImage(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get(":slug")
	fmt.Fprintf(w, "widgetImage %s\n", slug)
}
