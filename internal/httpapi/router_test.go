package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"http-server-projeto-korp/internal/httpapi"
)

func TestProjetoKorp(t *testing.T) {
	router, err := httpapi.NewRouter()
	if err != nil {
		t.Fatal(err)
	}

	// duas req ao mesmo router detectam um horário fixado na inicialização.
	for _, name := range []string{"primeira requisicao", "segunda requisicao"} {
		t.Run(name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/projeto-korp", nil)
			response := httptest.NewRecorder()

			before := time.Now().UTC()
			router.ServeHTTP(response, request)
			after := time.Now().UTC()

			if response.Code != http.StatusOK {
				t.Fatalf("status: recebido %d, esperado %d", response.Code, http.StatusOK)
			}
			if contentType := response.Header().Get("Content-Type"); contentType != "application/json; charset=utf-8" {
				t.Fatalf("Content-Type inesperado: %q", contentType)
			}

			var body map[string]string
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatalf("JSON inválido: %v", err)
			}
			if len(body) != 2 || body["nome"] != "Projeto Korp" {
				t.Fatalf("resposta inesperada: %v", body)
			}

			horario, err := time.Parse(time.RFC3339Nano, body["horario"])
			if err != nil {
				t.Fatalf("horario inválido: %v", err)
			}
			if horario.Location() != time.UTC {
				t.Fatalf("horario deve usar UTC com sufixo Z: %q", body["horario"])
			}
			if horario.Before(before) || horario.After(after) {
				t.Fatalf("horario %s fora do intervalo da requisição [%s, %s]", horario, before, after)
			}
		})
	}
}

func TestMetricsEndpoint(t *testing.T) {
	router, err := httpapi.NewRouter()
	if err != nil {
		t.Fatal(err)
	}
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/projeto-korp", nil))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("metrics status = %d", response.Code)
	}
	if !strings.Contains(response.Body.String(), `korp_http_requests_total{method="GET",route="/projeto-korp",status="200"} 1`) {
		t.Fatal("expected the API response counter in /metrics")
	}
	if strings.Contains(response.Body.String(), `route="/metrics"`) {
		t.Fatal("scrape requests must not be counted")
	}
}
