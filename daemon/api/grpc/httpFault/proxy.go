package httpfault

import (
	"errors"
	"io"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type faultTransport struct {
	wrapped   http.RoundTripper
	errorRate float32
	errorCode int
	delayMs   int32
	abortRate float32
}

func (t *faultTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.abortRate > 0 && rand.Float32() < t.abortRate {
		return nil, errors.New("injected abort")
	}
	if t.delayMs > 0 {
		time.Sleep(time.Duration(t.delayMs) * time.Millisecond)
	}
	if t.errorRate > 0 && rand.Float32() < t.errorRate {
		code := t.errorCode
		if code == 0 {
			code = http.StatusInternalServerError
		}
		return &http.Response{
			StatusCode: code,
			Status:     http.StatusText(code),
			Proto:      "HTTP/1.1",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(http.StatusText(code))),
		}, nil
	}
	return t.wrapped.RoundTrip(req)
}

// newFaultProxy builds a reverse proxy to targetURL with the given fault parameters.
func newFaultProxy(targetURL string, errorRate float32, errorCode int32, delayMs int32, abortRate float32) (http.Handler, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &faultTransport{
		wrapped:   http.DefaultTransport,
		errorRate: errorRate,
		errorCode: int(errorCode),
		delayMs:   delayMs,
		abortRate: abortRate,
	}
	return proxy, nil
}
