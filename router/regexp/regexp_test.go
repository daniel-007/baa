package regexp

import (
	"net/http"
	"testing"

	"github.com/go-baa/baa"
)

var b = baa.New()
var r = New(b)
var f = func(c *baa.Context) {}
var c = baa.NewContext(nil, nil, b)

func init() {
	b.SetDI("router", r)
}

func BenchmarkMatch1(b *testing.B) {
	router := loadBaaSingle("GET", "/user/:name", f)
	r, _ := http.NewRequest("GET", "/user/gordon", nil)
	benchRequest(b, router, r)
}

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

func loadBaaSingle(method, path string, h baa.HandlerFunc) http.Handler {
	b := baa.New()
	b.SetDI("router", New(b))
	switch method {
	case "GET":
		b.Get(path, h)
	case "POST":
		b.Post(path, h)
	case "PUT":
		b.Put(path, h)
	case "PATCH":
		b.Patch(path, h)
	case "DELETE":
		b.Delete(path, h)
	default:
		panic("Unknow HTTP method: " + method)
	}
	return b
}

func benchRequest(b *testing.B, router http.Handler, r *http.Request) {
	w := new(mockResponseWriter)
	u := r.URL
	rq := u.RawQuery
	r.RequestURI = u.RequestURI()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		u.RawQuery = rq
		router.ServeHTTP(w, r)
	}
}
