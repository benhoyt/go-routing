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
		{"POST", "/", 405, ""},

		{"GET", "/contact", 200, "contact\n"},
		{"POST", "/contact", 405, ""},
		{"GET", "/contact/", 404, ""},
		{"GET", "/contact/no", 404, ""},

		{"GET", "/api/widgets", 200, "apiGetWidgets\n"},
		{"GET", "/api/widgets/", 404, ""},

		{"POST", "/api/widgets", 200, "apiCreateWidget\n"},
		{"POST", "/api/widgets/", 404, ""},

		{"POST", "/api/widgets/foo", 200, "apiUpdateWidget foo\n"},
		{"POST", "/api/widgets/bar-baz", 200, "apiUpdateWidget bar-baz\n"},
		{"POST", "/api/widgets/foo/", 404, ""},
		{"GET", "/api/widgets/foo", 405, ""},

		{"POST", "/api/widgets/foo/parts", 200, "apiCreateWidgetPart foo\n"},
		{"POST", "/api/widgets/bar-baz/parts", 200, "apiCreateWidgetPart bar-baz\n"},
		{"POST", "/api/widgets/foo/parts/", 404, ""},
		{"GET", "/api/widgets/foo/parts", 405, ""},

		{"POST", "/api/widgets/foo/parts/1/update", 200, "apiUpdateWidgetPart foo 1\n"},
		{"POST", "/api/widgets/foo/parts/42/update", 200, "apiUpdateWidgetPart foo 42\n"},
		{"POST", "/api/widgets/foo/parts/bar/update", 404, ""},
		{"POST", "/api/widgets/bar-baz/parts/99/update", 200, "apiUpdateWidgetPart bar-baz 99\n"},
		{"GET", "/api/widgets/foo/parts/1/update", 405, ""},

		{"POST", "/api/widgets/foo/parts/1/delete", 200, "apiDeleteWidgetPart foo 1\n"},
		{"POST", "/api/widgets/foo/parts/42/delete", 200, "apiDeleteWidgetPart foo 42\n"},
		{"POST", "/api/widgets/foo/parts/bar/delete", 404, ""},
		{"POST", "/api/widgets/bar-baz/parts/99/delete", 200, "apiDeleteWidgetPart bar-baz 99\n"},
		{"GET", "/api/widgets/foo/parts/1/delete", 405, ""},

		{"GET", "/foo", 200, "widget foo\n"},
		{"GET", "/bar-baz", 200, "widget bar-baz\n"},
		{"GET", "/foo/", 404, ""},
		{"POST", "/foo", 405, ""},

		{"GET", "/foo/admin", 200, "widgetAdmin foo\n"},
		{"GET", "/bar-baz/admin", 200, "widgetAdmin bar-baz\n"},
		{"GET", "/foo/admin/", 404, ""},
		{"POST", "/foo/admin", 405, ""},

		{"POST", "/foo/image", 200, "widgetImage foo\n"},
		{"GET", "/foo/image", 405, ""},
		{"POST", "/bar-baz/image", 200, "widgetImage bar-baz\n"},
		{"POST", "/foo/image/", 404, ""},
		{"GET", "/foo/image", 405, ""},
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
					if test.status == 200 {
						body := recorder.Body.String()
						if body != test.body {
							t.Fatalf("expected body %q, got %q", test.body, body)
						}
					}
				})
			}
		})
	}
}

func BenchmarkRouters(b *testing.B) {
	tests := []struct {
		method string
		path   string
		status int
	}{
		{"GET", "/", 200},
		{"GET", "/api/widgets", 200},
		{"POST", "/api/widgets/foo", 200},
		{"POST", "/api/widgets/foo/parts/1/update", 200},
		{"GET", "/foo", 200},
	}
	for _, name := range routerNames {
		router := routers[name]
		b.Run(name, func(b *testing.B) {
			for _, test := range tests {
				path := strings.ReplaceAll(test.path, "/", "_")
				b.Run(test.method+path, func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						recorder := httptest.NewRecorder()
						request, err := http.NewRequest(test.method, test.path, &bytes.Buffer{})
						if err != nil {
							b.Fatal(err)
						}
						router.ServeHTTP(recorder, request)
						if recorder.Code != test.status {
							b.Fatalf("expected status %d, got %d", test.status, recorder.Code)
						}
					}
				})
			}
		})
	}
}
