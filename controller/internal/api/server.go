package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	v1 "obzev0/controller/api/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StartAPIServer starts the dashboard HTTP server in a background goroutine.
// It serves the embedded React SPA on / and a REST API on /api/.
func StartAPIServer(mgr ctrl.Manager, port int) {
	mux := http.NewServeMux()

	c := mgr.GetClient()

	mux.HandleFunc("/api/experiments", func(w http.ResponseWriter, r *http.Request) {
		listExperiments(w, r, c)
	})
	mux.HandleFunc("/api/pods", func(w http.ResponseWriter, r *http.Request) {
		listPods(w, r, c)
	})
	mux.HandleFunc("/api/experiments/", func(w http.ResponseWriter, r *http.Request) {
		experimentAction(w, r, c)
	})
	mux.HandleFunc("/api/reports", func(w http.ResponseWriter, r *http.Request) {
		listReports(w, r, c)
	})
	mux.HandleFunc("/api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Serve embedded React SPA — falls back gracefully when no UI dist is present.
	uiHandler := buildUIHandler()
	mux.Handle("/", uiHandler)

	go func() {
		addr := fmt.Sprintf(":%d", port)
		log.Printf("Dashboard listening on %s", addr)
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Printf("Dashboard server error: %v", err)
		}
	}()
}

func jsonResponse(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(v)
}

// listExperiments returns all Obzev0Resource CRs across all namespaces.
func listExperiments(w http.ResponseWriter, r *http.Request, c client.Client) {
	var list v1.Obzev0ResourceList
	if err := c.List(context.Background(), &list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type experimentSummary struct {
		Name       string `json:"name"`
		Namespace  string `json:"namespace"`
		Schedule   string `json:"schedule,omitempty"`
		StepCount  int    `json:"stepCount,omitempty"`
		HasRollback bool  `json:"hasRollback"`
		BlastRadius v1.BlastRadius `json:"blastRadius,omitempty"`
	}

	out := make([]experimentSummary, len(list.Items))
	for i, item := range list.Items {
		out[i] = experimentSummary{
			Name:        item.Name,
			Namespace:   item.Namespace,
			Schedule:    item.Spec.Schedule,
			StepCount:   len(item.Spec.Steps),
			HasRollback: item.Spec.RollbackPolicy.MaxDurationSeconds > 0 || item.Spec.RollbackPolicy.PrometheusURL != "",
			BlastRadius: item.Spec.BlastRadius,
		}
	}
	jsonResponse(w, out)
}

// listPods returns all daemon pods (label app=grpc-server).
func listPods(w http.ResponseWriter, r *http.Request, c client.Client) {
	var podList corev1.PodList
	if err := c.List(context.Background(), &podList, &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set{"app": "grpc-server"}),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type podSummary struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
		Node      string `json:"node"`
		Phase     string `json:"phase"`
		IP        string `json:"ip"`
	}

	out := make([]podSummary, len(podList.Items))
	for i, pod := range podList.Items {
		out[i] = podSummary{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Node:      pod.Spec.NodeName,
			Phase:     string(pod.Status.Phase),
			IP:        pod.Status.PodIP,
		}
	}
	jsonResponse(w, out)
}

// experimentAction handles DELETE /api/experiments/{namespace}/{name} to stop an experiment.
func experimentAction(w http.ResponseWriter, r *http.Request, c client.Client) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Path: /api/experiments/{namespace}/{name}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/experiments/"), "/")
	if len(parts) != 2 {
		http.Error(w, "path must be /api/experiments/{namespace}/{name}", http.StatusBadRequest)
		return
	}
	ns, name := parts[0], parts[1]

	var cr v1.Obzev0Resource
	cr.Name = name
	cr.Namespace = ns
	if err := c.Delete(context.Background(), &cr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// listReports returns ConfigMaps with label obzev0.io/report=true.
func listReports(w http.ResponseWriter, r *http.Request, c client.Client) {
	var cmList corev1.ConfigMapList
	if err := c.List(context.Background(), &cmList, &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set{"obzev0.io/report": "true"}),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type reportSummary struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
		Report    string `json:"report"`
	}
	out := make([]reportSummary, len(cmList.Items))
	for i, cm := range cmList.Items {
		out[i] = reportSummary{
			Name:      cm.Name,
			Namespace: cm.Namespace,
			Report:    cm.Data["report.json"],
		}
	}
	jsonResponse(w, out)
}

// buildUIHandler returns a handler that serves the embedded React dist,
// or a placeholder page if no dist is present.
func buildUIHandler() http.Handler {
	sub, err := fs.Sub(uiFS, "dist")
	if err != nil {
		// No UI dist available — serve a simple status page.
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<!DOCTYPE html><html><body>
<h2>obzev0 dashboard</h2>
<p>UI not built. Run <code>cd ui && npm run build</code> and rebuild the controller image.</p>
<p><a href="/api/experiments">GET /api/experiments</a></p>
<p><a href="/api/pods">GET /api/pods</a></p>
<p><a href="/api/reports">GET /api/reports</a></p>
</body></html>`))
		})
	}
	fileServer := http.FileServer(http.FS(sub))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// SPA: serve index.html for unknown paths so React Router works.
		if _, err := sub.Open(strings.TrimPrefix(r.URL.Path, "/")); err != nil {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
}
