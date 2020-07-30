// Test various ways to do HTTP method+path routing in Go

// Each router handles the 11 URLs below:
//
// GET  /										# home
// GET  /contact								# contact
// GET  /api/widgets							# apiGetWidgets
// POST /api/widgets                           	# apiCreateWidget
// POST /api/widgets/:slug                     	# apiUpdateWidget
// POST /api/widgets/:slug/parts               	# apiCreateWidgetPart
// POST /api/widgets/:slug/parts/:id/update    	# apiUpdateWidgetPart
// POST /api/widgets/:slug/parts/:id/delete    	# apiDeleteWidgetPart
// GET  /:slug									# widget
// GET  /:slug/admin                         	# widgetAdmin
// POST /:slug/image							# widgetImage

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/benhoyt/go-routing/chi"
	"github.com/benhoyt/go-routing/gorilla"
	"github.com/benhoyt/go-routing/match"
	"github.com/benhoyt/go-routing/pat"
	"github.com/benhoyt/go-routing/reswitch"
	"github.com/benhoyt/go-routing/retable"
	"github.com/benhoyt/go-routing/shiftpath"
	"github.com/benhoyt/go-routing/split"
)

const port = 8080

func main() {
	if len(os.Args) < 2 || routers[os.Args[1]] == nil {
		fmt.Fprintf(os.Stderr, "usage: go-routing router\n\n")
		fmt.Fprintf(os.Stderr, "router is one of: %s\n", strings.Join(routerNames, ", "))
		os.Exit(1)
	}
	routerName := os.Args[1]
	router := routers[routerName]

	fmt.Printf("listening on port %d using %s router\n", port, routerName)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

var routers = map[string]http.Handler{
	"chi":       chi.Serve,
	"gorilla":   gorilla.Serve,
	"match":     http.HandlerFunc(match.Serve),
	"pat":       pat.Serve,
	"reswitch":  http.HandlerFunc(reswitch.Serve),
	"retable":   http.HandlerFunc(retable.Serve),
	"shiftpath": http.HandlerFunc(shiftpath.Serve),
	"split":     http.HandlerFunc(split.Serve),
}

var routerNames = func() []string {
	routerNames := []string{}
	for k := range routers {
		routerNames = append(routerNames, k)
	}
	sort.Strings(routerNames)
	return routerNames
}()
