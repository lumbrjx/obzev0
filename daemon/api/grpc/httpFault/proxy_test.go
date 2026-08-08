package httpfault

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFaultProxy_NoFault(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}))
	defer backend.Close()

	proxy, err := newFaultProxy(backend.URL, 0, 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	front := httptest.NewServer(proxy)
	defer front.Close()

	resp, err := http.Get(front.URL + "/test")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "ok" {
		t.Fatalf("expected body 'ok', got %q", string(body))
	}
}

func TestFaultProxy_ErrorRate100(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}))
	defer backend.Close()

	proxy, err := newFaultProxy(backend.URL, 1.0, 503, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	front := httptest.NewServer(proxy)
	defer front.Close()

	for i := 0; i < 10; i++ {
		resp, err := http.Get(front.URL + "/test")
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != 503 {
			t.Fatalf("request %d: expected 503, got %d", i, resp.StatusCode)
		}
	}
}

func TestFaultProxy_AbortRate100(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer backend.Close()

	proxy, err := newFaultProxy(backend.URL, 0, 0, 0, 1.0)
	if err != nil {
		t.Fatal(err)
	}
	front := httptest.NewServer(proxy)
	defer front.Close()

	for i := 0; i < 5; i++ {
		resp, err := http.Get(front.URL + "/test")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode < 500 {
				t.Fatalf("request %d: expected an error or 5xx, got %d", i, resp.StatusCode)
			}
		}
	}
}

func TestFaultProxy_Delay(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer backend.Close()

	proxy, err := newFaultProxy(backend.URL, 0, 0, 100, 0)
	if err != nil {
		t.Fatal(err)
	}
	front := httptest.NewServer(proxy)
	defer front.Close()

	start := time.Now()
	resp, err := http.Get(front.URL + "/test")
	elapsed := time.Since(start)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if elapsed < 80*time.Millisecond {
		t.Fatalf("expected at least ~100ms delay, got %s", elapsed)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestFaultProxy_PartialErrorRate(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer backend.Close()

	proxy, err := newFaultProxy(backend.URL, 0.5, 500, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	front := httptest.NewServer(proxy)
	defer front.Close()

	errors := 0
	total := 200
	for i := 0; i < total; i++ {
		resp, err := http.Get(front.URL + "/test")
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode == 500 {
			errors++
		}
		resp.Body.Close()
	}

	errorRate := float64(errors) / float64(total)
	if errorRate < 0.2 || errorRate > 0.8 {
		t.Fatalf("expected ~50%% error rate, got %.1f%% (%d/%d)", errorRate*100, errors, total)
	}
}

func TestFaultProxy_DefaultErrorCode(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer backend.Close()

	proxy, err := newFaultProxy(backend.URL, 1.0, 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	front := httptest.NewServer(proxy)
	defer front.Close()

	resp, err := http.Get(front.URL + "/test")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 500 {
		t.Fatalf("expected default 500, got %d", resp.StatusCode)
	}
}
