package limiter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"rate-limiter/config"
	"rate-limiter/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestMiddleware() (*gin.Engine, *RateLimiterService) {
	gin.SetMode(gin.TestMode)

	config.InitLogger()

	// Criar armazenamento em memória para os testes
	memStorage := storage.NewMemoryStorage()

	// Criar serviço de Rate Limiter
	rateLimiter := NewRateLimiterService(memStorage, RateLimiterConfig{
		RateLimitPerIP:        3, // Limite baixo para facilitar os testes
		RateLimitPerToken:     5,
		DefaultBlockTimeIP:    120,
		DefaultBlockTimeToken: 300,
	}, config.Logger)

	// Criar servidor Gin de teste
	router := gin.New()
	router.Use(RateLimiterMiddleware(rateLimiter))

	// Rota de teste
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Requisição permitida"})
	})

	return router, rateLimiter
}

func TestMiddleware_IPBlocking(t *testing.T) {
	router, _ := setupTestMiddleware()

	// Criar requisição simulada
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	// Fazer 3 requisições válidas (dentro do limite)
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// 4ª requisição deve ser bloqueada
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestMiddleware_TokenBlocking(t *testing.T) {
	router, _ := setupTestMiddleware()

	// Criar requisição simulada com Token
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("API_KEY", "test_token")

	// Fazer 5 requisições válidas (dentro do limite)
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// 6ª requisição deve ser bloqueada
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestMiddleware_TokenTakesPriorityOverIP(t *testing.T) {
	router, _ := setupTestMiddleware()

	// Criar requisição simulada com Token e IP bloqueado
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	// **Bloquear o IP** com 3 requisições
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder() // Resetando para cada requisição
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Esperado 200 OK na requisição %d", i+1)
	}

	// **Fazer a 4ª requisição que deve ser bloqueada**
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "Esperado 429 Too Many Requests quando o IP atingir o limite")

	// **Agora testamos se o Token tem prioridade**
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	req2.Header.Set("API_KEY", "token_priority")

	// **As próximas 5 requisições com Token devem ser permitidas**
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder() // Resetando para cada requisição
		router.ServeHTTP(w, req2)
		assert.Equal(t, http.StatusOK, w.Code, "Esperado 200 OK para requisição com Token na tentativa %d", i+1)
	}

	// **Agora o Token deve atingir seu limite e ser bloqueado**
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req2)
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "Esperado 429 Too Many Requests quando o Token atingir o limite")
}

func TestMiddleware_ResponseBody(t *testing.T) {
	router, _ := setupTestMiddleware()

	// Criar requisição simulada
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	// Fazer 3 requisições válidas
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// 4ª requisição deve ser bloqueada
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Contains(t, w.Body.String(), "You have reached the maximum number of requests or actions allowed within a certain time frame")
}
