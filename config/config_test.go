package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Simulando variáveis de ambiente para o teste
	os.Setenv("RATE_LIMIT_PER_IP", "10")
	os.Setenv("RATE_LIMIT_PER_TOKEN", "200")
	os.Setenv("DEFAULT_BLOCK_TIME_IP", "180")
	os.Setenv("DEFAULT_BLOCK_TIME_TOKEN", "360")
	os.Setenv("BLOCK_TIME_PER_IP", "192.168.1.1=120;192.168.1.2=600")
	os.Setenv("BLOCK_TIME_PER_TOKEN", "token123=180;tokenABC=900")
	os.Setenv("REDIS_ADDR", "redis-test:6379")
	os.Setenv("SERVER_PORT", "9090")

	// Carregar configurações
	LoadConfig()

	// Testando valores carregados corretamente
	assert.Equal(t, 10, Cfg.RateLimitPerIP)
	assert.Equal(t, 200, Cfg.RateLimitPerToken)
	assert.Equal(t, 180, Cfg.DefaultBlockTimeIP)
	assert.Equal(t, 360, Cfg.DefaultBlockTimeToken)
	assert.Equal(t, "redis-test:6379", Cfg.RedisAddr)
	assert.Equal(t, "9090", Cfg.ServerPort)

	// Testando tempos de bloqueio individuais
	assert.Equal(t, 120, Cfg.BlockTimePerIP["192.168.1.1"])
	assert.Equal(t, 600, Cfg.BlockTimePerIP["192.168.1.2"])
	assert.Equal(t, 180, Cfg.BlockTimePerToken["token123"])
	assert.Equal(t, 900, Cfg.BlockTimePerToken["tokenABC"])

	// Testando valores padrão para IPs e Tokens não especificados
	assert.Equal(t, 180, GetBlockTimeForIP("192.168.1.3")) // Deve retornar DEFAULT_BLOCK_TIME_IP
	assert.Equal(t, 360, GetBlockTimeForToken("tokenXYZ")) // Deve retornar DEFAULT_BLOCK_TIME_TOKEN
}

// GetBlockTimeForIP verifica se um IP tem um tempo de bloqueio específico ou retorna o padrão
func GetBlockTimeForIP(ip string) int {
	if time, exists := Cfg.BlockTimePerIP[ip]; exists {
		return time
	}
	return Cfg.DefaultBlockTimeIP
}

// GetBlockTimeForToken verifica se um Token tem um tempo de bloqueio específico ou retorna o padrão
func GetBlockTimeForToken(token string) int {
	if time, exists := Cfg.BlockTimePerToken[token]; exists {
		return time
	}
	return Cfg.DefaultBlockTimeToken
}
