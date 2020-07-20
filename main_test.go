// Test the routers

package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouters(t *testing.T) {
	tests := []struct {
		method string
		path   string
		status int
		body   string
	}{
		{"GET", "/", 200, "home\n"},
		{"HEAD", "/", 200, "home\n"},
		{"POST", "/", 405, "405 method not allowed\n"},

		{"GET", "/contact", 200, "contact\n"},
		{"HEAD", "/contact", 200, "contact\n"},
		{"POST", "/contact", 405, "405 method not allowed\n"},
		{"GET", "/contact/", 404, "404 page not found\n"},
		{"GET", "/contact/no", 404, "404 page not found\n"},

		{"GET", "/api/widgets", 200, "apiGetWidgets\n"},
		{"HEAD", "/api/widgets", 200, "apiGetWidgets\n"},
		{"GET", "/api/widgets/", 404, "404 page not found\n"},

		{"POST", "/api/widgets", 200, "apiCreateWidget\n"},
		{"POST", "/api/widgets/", 404, "404 page not found\n"},

		{"POST", "/api/widgets/foo", 200, "apiUpdateWidget foo\n"},
		{"POST", "/api/widgets/bar-baz", 200, "apiUpdateWidget bar-baz\n"},
		{"POST", "/api/widgets/foo/", 404, "404 page not found\n"},
		{"GET", "/api/widgets/foo", 405, "405 method not allowed\n"},

		{"POST", "/api/widgets/foo/parts", 200, "apiCreateWidgetPart foo\n"},
		{"POST", "/api/widgets/bar-baz/parts", 200, "apiCreateWidgetPart bar-baz\n"},
		{"POST", "/api/widgets/foo/parts/", 404, "404 page not found\n"},
		{"GET", "/api/widgets/foo/parts", 405, "405 method not allowed\n"},

		{"POST", "/api/widgets/foo/parts/1/update", 200, "apiUpdateWidgetPart foo 1\n"},
		{"POST", "/api/widgets/foo/parts/42/update", 200, "apiUpdateWidgetPart foo 42\n"},
		{"POST", "/api/widgets/bar-baz/parts/99/update", 200, "apiUpdateWidgetPart bar-baz 99\n"},
		{"GET", "/api/widgets/foo/parts/1/update", 405, "405 method not allowed\n"},

		{"POST", "/api/widgets/foo/parts/1/delete", 200, "apiDeleteWidgetPart foo 1\n"},
		{"POST", "/api/widgets/foo/parts/42/delete", 200, "apiDeleteWidgetPart foo 42\n"},
		{"POST", "/api/widgets/bar-baz/parts/99/delete", 200, "apiDeleteWidgetPart bar-baz 99\n"},
		{"GET", "/api/widgets/foo/parts/1/delete", 405, "405 method not allowed\n"},

		{"GET", "/foo", 200, "widget foo\n"},
		{"HEAD", "/foo", 200, "widget foo\n"},
		{"GET", "/bar-baz", 200, "widget bar-baz\n"},
		{"GET", "/foo/", 404, "404 page not found\n"},
		{"POST", "/foo", 405, "405 method not allowed\n"},

		{"GET", "/foo/admin", 200, "widgetAdmin foo\n"},
		{"HEAD", "/foo/admin", 200, "widgetAdmin foo\n"},
		{"GET", "/bar-baz/admin", 200, "widgetAdmin bar-baz\n"},
		{"GET", "/foo/admin/", 404, "404 page not found\n"},
		{"POST", "/foo/admin", 405, "405 method not allowed\n"},

		{"POST", "/foo/image", 200, "widgetImage foo\n"},
		{"GET", "/foo/image", 405, "405 method not allowed\n"},
		{"POST", "/bar-baz/image", 200, "widgetImage bar-baz\n"},
		{"POST", "/foo/image/", 404, "404 page not found\n"},
		{"GET", "/foo/image", 405, "405 method not allowed\n"},
	}
	for _, name := range routerNames {
		router := routers[name]
		t.Run(name, func(t *testing.T) {
			for _, test := range tests {
				path := strings.ReplaceAll(test.path, "/", "_")
				t.Run(test.method+path, func(t *testing.T) {
					recorder := httptest.NewRecorder()
					request, err := http.NewRequest(test.method, test.path, &bytes.Buffer{})
					if err != nil {
						t.Fatal(err)
					}
					router.ServeHTTP(recorder, request)
					if recorder.Code != test.status {
						t.Fatalf("expected status %d, got %d", test.status, recorder.Code)
					}
					body := recorder.Body.String()
					if body != test.body {
						t.Fatalf("expected body %q, got %q", test.body, body)
					}
				})
			}
		})
	}
}
