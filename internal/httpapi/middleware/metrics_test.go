package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"http-server-projeto-korp/internal/httpapi/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func TestMetricsRecordsResponses(t *testing.T) {
	for _, tc := range []struct {
		name   string
		path   string
		status int
		route  string
	}{
		{"success", "/ok", http.StatusOK, "/ok"},
		{"not found", "/missing", http.StatusNotFound, "unmatched"},
		{"recovered panic", "/panic", http.StatusInternalServerError, "/panic"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			registry := prometheus.NewRegistry()
			router := gin.New()
			router.Use(middleware.Metrics(registry), gin.Recovery())
			router.GET("/ok", func(c *gin.Context) { c.Status(http.StatusOK) })
			router.GET("/panic", func(c *gin.Context) { panic("test") })

			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d", response.Code, tc.status)
			}

			families, err := registry.Gather()
			if err != nil {
				t.Fatal(err)
			}
			if len(families) != 2 {
				t.Fatalf("metric families = %d, want 2", len(families))
			}
			for _, family := range families {
				if len(family.Metric) != 1 {
					t.Fatalf("unexpected series in %s", family.GetName())
				}
				metric := family.Metric[0]
				labels := make(map[string]string)
				for _, label := range metric.Label {
					labels[label.GetName()] = label.GetValue()
				}
				if labels["route"] != tc.route || labels["method"] != http.MethodGet {
					t.Fatalf("unexpected labels: %v", labels)
				}
				switch family.GetName() {
				case "korp_http_requests_total":
					if metric.GetCounter().GetValue() != 1 || labels["status"] != strconv.Itoa(tc.status) {
						t.Fatalf("unexpected response counter: %v", metric)
					}
				case "korp_http_request_duration_seconds":
					if metric.GetHistogram().GetSampleCount() != 1 || metric.GetHistogram().GetSampleSum() < 0 {
						t.Fatalf("unexpected duration histogram: %v", metric)
					}
				}
			}
		})
	}
}

func TestMetricsExcludesScrapes(t *testing.T) {
	registry := prometheus.NewRegistry()
	router := gin.New()
	router.Use(middleware.Metrics(registry))
	router.GET("/metrics", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/metrics", nil))
	families, err := registry.Gather()
	if err != nil {
		t.Fatal(err)
	}
	if len(families) != 0 {
		t.Fatal("metrics scrapes must not count as application traffic")
	}
}
