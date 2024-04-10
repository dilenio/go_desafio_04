package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dilenio/desafio04/limiter"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/stretchr/testify/assert"
)

func TestEndpointWithValidToken(t *testing.T) {
	router := chi.NewRouter()
	rateLimiter := limiter.NewRateLimiter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(rateLimiter.LimitHandler)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	// Faz uma requisição GET ao servidor de teste com API_KEY válida
	for i := 0; i < 10; i++ {
		req, err := http.NewRequest("GET", ts.URL, nil)
		if err != nil {
			t.Fatalf("Erro ao criar requisição: %v", err)
			return
		}

		req.Header.Add("API_KEY", "abc123456")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Erro ao fazer a requisição GET: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Erro ao ler o corpo da resposta: %v", err)
		}

		if i <= 4 {
			assert.Equal(t, http.StatusOK, resp.StatusCode, "O status da resposta deveria ser 200 OK")
			assert.Equal(t, "Hello World!", string(body), "O corpo da resposta deveria ser 'Hello World!'")
		} else {
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode, "O status da resposta deveria ser 429 Too Many Requests")
			assert.Equal(t, "you have reached the maximum number of requests or actions allowed within a certain time frame\n", string(body), "O corpo da resposta deveria ser 'you have reached the maximum number of requests or actions allowed within a certain time frame'")
		}
	}

	// Faz uma requisição GET ao servidor de teste com API_KEY inválida
	for i := 0; i < 10; i++ {
		req, err := http.NewRequest("GET", ts.URL, nil)
		if err != nil {
			t.Fatalf("Erro ao criar requisição: %v", err)
			return
		}

		req.Header.Add("API_KEY", "xyz987654")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Erro ao fazer a requisição GET: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Erro ao ler o corpo da resposta: %v", err)
		}

		if i <= 1 {
			assert.Equal(t, http.StatusOK, resp.StatusCode, "O status da resposta deveria ser 200 OK")
			assert.Equal(t, "Hello World!", string(body), "O corpo da resposta deveria ser 'Hello World!'")
		} else {
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode, "O status da resposta deveria ser 429 Too Many Requests")
			assert.Equal(t, "you have reached the maximum number of requests or actions allowed within a certain time frame\n", string(body), "O corpo da resposta deveria ser 'you have reached the maximum number of requests or actions allowed within a certain time frame'")
		}
	}

	// AGUARDA 6 SEGUNDOS PARA FAZER MAIS REQUISIÇÕES
	t.Log("Aguardando 6 segundos para fazer mais requisições...")
	time.Sleep(6 * time.Second)

	// Faz uma requisição GET ao servidor de teste apenas com IP
	for i := 0; i < 10; i++ {
		req, err := http.NewRequest("GET", ts.URL, nil)
		if err != nil {
			t.Fatalf("Erro ao criar requisição: %v", err)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Erro ao fazer a requisição GET: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Erro ao ler o corpo da resposta: %v", err)
		}

		if i <= 1 {
			assert.Equal(t, http.StatusOK, resp.StatusCode, "O status da resposta deveria ser 200 OK")
			assert.Equal(t, "Hello World!", string(body), "O corpo da resposta deveria ser 'Hello World!'")
		} else {
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode, "O status da resposta deveria ser 429 Too Many Requests")
			assert.Equal(t, "you have reached the maximum number of requests or actions allowed within a certain time frame\n", string(body), "O corpo da resposta deveria ser 'you have reached the maximum number of requests or actions allowed within a certain time frame'")
		}
	}
}
