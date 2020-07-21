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

// TODO: Axel Wagner's version

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/benhoyt/go-routing/routechi"
	"github.com/benhoyt/go-routing/routegorilla"
	"github.com/benhoyt/go-routing/routematch"
	"github.com/benhoyt/go-routing/routepat"
	"github.com/benhoyt/go-routing/routeregexswitch"
	"github.com/benhoyt/go-routing/routeregextable"
	"github.com/benhoyt/go-routing/routesplit"
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

	http.Handle("/", router)
	fmt.Printf("listening on port %d using %s router\n", port, routerName)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

var routers = map[string]http.Handler{
	"chi":          routechi.Route,
	"gorilla":      routegorilla.Route,
	"match":        http.HandlerFunc(routematch.Route),
	"pat":          routepat.Route,
	"regex-switch": http.HandlerFunc(routeregexswitch.Route),
	"regex-table":  http.HandlerFunc(routeregextable.Route),
	"split-flat":   http.HandlerFunc(routesplit.RouteFlat),
	"split-nested": http.HandlerFunc(routesplit.RouteNested),
}

var routerNames = func() []string {
	routerNames := []string{}
	for k := range routers {
		routerNames = append(routerNames, k)
	}
	sort.Strings(routerNames)
	return routerNames
}()
