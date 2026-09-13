package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
