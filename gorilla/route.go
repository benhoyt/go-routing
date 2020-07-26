// Routing based on the gorilla/mux router

package gorilla

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var Serve http.Handler

func init() {
	r := mux.NewRouter()

	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/contact", contact).Methods("GET")
	r.HandleFunc("/api/widgets", apiGetWidgets).Methods("GET")
	r.HandleFunc("/api/widgets", apiCreateWidget).Methods("POST")
	r.HandleFunc("/api/widgets/{slug}", apiUpdateWidget).Methods("POST")
	r.HandleFunc("/api/widgets/{slug}/parts", apiCreateWidgetPart).Methods("POST")
	r.HandleFunc("/api/widgets/{slug}/parts/{id:[0-9]+}/update", apiUpdateWidgetPart).Methods("POST")
	r.HandleFunc("/api/widgets/{slug}/parts/{id:[0-9]+}/delete", apiDeleteWidgetPart).Methods("POST")
	r.HandleFunc("/{slug}", widgetGet).Methods("GET")
	r.HandleFunc("/{slug}/admin", widgetAdmin).Methods("GET")
	r.HandleFunc("/{slug}/image", widgetImage).Methods("POST")

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
	slug := mux.Vars(r)["slug"]
	fmt.Fprintf(w, "apiUpdateWidget %s\n", slug)
}

func apiCreateWidgetPart(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	fmt.Fprintf(w, "apiCreateWidgetPart %s\n", slug)
}

func apiUpdateWidgetPart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	id, _ := strconv.Atoi(vars["id"])
	fmt.Fprintf(w, "apiUpdateWidgetPart %s %d\n", slug, id)
}

func apiDeleteWidgetPart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	id, _ := strconv.Atoi(vars["id"])
	fmt.Fprintf(w, "apiDeleteWidgetPart %s %d\n", slug, id)
}

func widgetGet(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	fmt.Fprintf(w, "widget %s\n", slug)
}

func widgetAdmin(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	fmt.Fprintf(w, "widgetAdmin %s\n", slug)
}

func widgetImage(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	fmt.Fprintf(w, "widgetImage %s\n", slug)
}
