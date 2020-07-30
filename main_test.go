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
		{"POST", "/api/widgets/foo/zarts", 404, ""},

		{"POST", "/api/widgets/foo/parts/1/update", 200, "apiUpdateWidgetPart foo 1\n"},
		{"POST", "/api/widgets/foo/parts/1/update/no", 404, ""},
		{"POST", "/api/widgets/foo/parts/42/update", 200, "apiUpdateWidgetPart foo 42\n"},
		{"POST", "/api/widgets/foo/parts/bar/update", 404, ""},
		{"POST", "/api/widgets/bar-baz/parts/99/update", 200, "apiUpdateWidgetPart bar-baz 99\n"},
		{"GET", "/api/widgets/foo/parts/1/update", 405, ""},

		{"POST", "/api/widgets/foo/parts/1/delete", 200, "apiDeleteWidgetPart foo 1\n"},
		{"POST", "/api/widgets/foo/parts/1/delete/no", 404, ""},
		{"POST", "/api/widgets/foo/parts/42/delete", 200, "apiDeleteWidgetPart foo 42\n"},
		{"POST", "/api/widgets/foo/parts/bar/delete", 404, ""},
		{"POST", "/api/widgets/bar-baz/parts/99/delete", 200, "apiDeleteWidgetPart bar-baz 99\n"},
		{"GET", "/api/widgets/foo/parts/1/delete", 405, ""},
		{"POST", "/api/widgets/foo/parts/1/no", 404, ""},

		{"GET", "/foo", 200, "widget foo\n"},
		{"GET", "/bar-baz", 200, "widget bar-baz\n"},
		{"GET", "/foo/", 404, ""},
		{"POST", "/foo", 405, ""},

		{"GET", "/foo/admin", 200, "widgetAdmin foo\n"},
		{"GET", "/bar-baz/admin", 200, "widgetAdmin bar-baz\n"},
		{"GET", "/foo/admin/", 404, ""},
		{"GET", "/foo/admin/no", 404, ""},
		{"POST", "/foo/admin", 405, ""},

		{"POST", "/foo/image", 200, "widgetImage foo\n"},
		{"POST", "/foo/image/no", 404, ""},
		{"GET", "/foo/image", 405, ""},
		{"POST", "/bar-baz/image", 200, "widgetImage bar-baz\n"},
		{"POST", "/foo/image/", 404, ""},
		{"GET", "/foo/image", 405, ""},
		{"GET", "/foo/no", 404, ""},
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
	method := "POST"
	path := "/api/widgets/foo/parts/1/update"

	// Could use httptest.ResponseRecorder, but that's slow-ish
	responseWriter := &noopResponseWriter{}

	names := make([]string, len(routerNames))
	copy(names, routerNames)
	names = append(names, "noop")

	for _, name := range names {
		router := routers[name]
		if router == nil {
			router = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		}
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				request, err := http.NewRequest(method, path, &bytes.Buffer{})
				if err != nil {
					b.Fatal(err)
				}
				router.ServeHTTP(responseWriter, request)
			}
		})
	}
}

type noopResponseWriter struct{}

func (r *noopResponseWriter) Header() http.Header {
	return nil
}

func (r *noopResponseWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (r *noopResponseWriter) WriteHeader(statusCode int) {
}
