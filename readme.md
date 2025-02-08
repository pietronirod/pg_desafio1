# 🚀 Rate Limiter em Go

Este projeto implementa um **Rate Limiter** em Go para limitar requisições por **IP** e **Token**.  
Utiliza **Redis** para armazenamento e pode ser configurado via `.env`.

## 📌 Funcionalidades

- Limitação de requisições por **IP** ou **Token**.
- Bloqueio temporário após exceder o limite.
- Middleware para integração com **Gin**.
- Suporte a **Redis** para escalabilidade.
- Testes unitários e de integração.

---

## 📦 1. Configuração do Ambiente

### **1️⃣ Pré-requisitos**

- **Go** 1.23
- **Docker e Docker Compose** (para Redis)

### **2️⃣ Clonar o repositório**

```sh
git clone https://github.com/pietronirod/pg_desafio1.git
cd rate-limiter
```

### **3️⃣ Criar o arquivo `.env`**

Crie um arquivo `.env` na raiz do projeto com as configurações:

```ini
# Limites de requisições
RATE_LIMIT_PER_IP=5
RATE_LIMIT_PER_TOKEN=100

# Tempo de bloqueio em segundos (ex: 300 = 5 minutos)
DEFAULT_BLOCK_TIME_IP=300
DEFAULT_BLOCK_TIME_TOKEN=600

# Lista de bloqueios individuais (Formato IP=tempo_em_segundos, separado por ;)
BLOCK_TIME_PER_IP=192.168.1.1=120;192.168.1.2=600;192.168.1.3=900

# Lista de bloqueios por token (Formato token=tempo_em_segundos, separado por ;)
BLOCK_TIME_PER_TOKEN=token1=120;token2=600;token3=900

# Configuração do Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0

# Configuração do servidor
SERVER_PORT=8080

# Configuração de logging
LOG_LEVEL=info
```

### **4️⃣ Subir o Redis com Docker**

```sh
docker-compose up --build -d
```

Para verificar os logs

```sh
docker-compose logs -f
```

---

## 🚀 2. Rodando o Projeto

### **Iniciar o servidor**

```sh
go run cmd/main.go
```

O servidor será iniciado na porta `8080` (ou conforme definido no `.env`).

### **Testando o Rate Limiter**

Abra um terminal e execute:

```sh
curl -i http://localhost:8080/
```

Após **10 requisições rápidas**, o Rate Limiter retornará:

```json
{
  "message": "You have reached the maximum number of requests or actions allowed within a certain time frame",
  "retry_in": 300
}
```

Se quiser testar **com Token**, adicione um cabeçalho:

```sh
curl -i -H "API_KEY: meu_token" http://localhost:8080/
```

---

## 🛠️ 3. Rodando os Testes

### **Rodar todos os testes**

```sh
go test ./... -v
```

### **Rodar testes unitários**

```sh
go test ./internal/limiter -v
```

### **Rodar testes de integração com Redis**

```sh
go test ./internal/storage/redis_integration_test.go -v
```

---

## 📌 4. Estrutura do Projeto

```bash
/rate-limiter-go
├── cmd/
│   ├── main.go  # Entrada do servidor
│
├── config/
│   ├── config.go  # Carregamento de configurações
│   ├── logger.go  # Configuração do logging estruturado
│
├── internal/
│   ├── limiter/
│   │   ├── model.go       # Estruturas de dados
│   │   ├── service.go     # Lógica do Rate Limiter
│   │   ├── middleware.go  # Middleware para Gin
│   │   ├── load_test.go   # Testes de carga
│   │
│   ├── storage/
│   │   ├── storage.go      # Interface de persistência
│   │   ├── redis.go        # Implementação do Redis
│   │   ├── memory.go       # Implementação em memória
│   │   ├── redis_integration_test.go  # Testes de integração com Redis
│
├── .env  # Configuração de ambiente
├── docker-compose.yml  # Configuração do Redis
├── go.mod  # Dependências do projeto
├── README.md  # Documentação
```

---

## 🔥 5. Melhorias Futuras

- [ ] Implementar um **painel de monitoramento** para visualizar as requisições bloqueadas.
- [ ] Suporte a **limitação por rota** além de IP e Token.
- [ ] Melhorias na **persistência** para suportar outros bancos além do Redis.

---

## 🐝 Licença

Este projeto é distribuído sob a licença **MIT**.

---

## 🚀 Contribuições

Sinta-se à vontade para contribuir! Sugestões e melhorias são sempre bem-vindas.  
