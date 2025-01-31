## Desafio #07 - Rate Limit

O objetivo deste desafio é criar um rate limiter em Go que possa ser utilizado para controlar o tráfego de requisições para um serviço web. .

### Requisitos:

O rate limiter deve ser capaz de limitar o número de requisições com base em dois critérios:
 
1. **Endereço IP:** O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.
2. **Token de Acesso:** O rate limiter deve também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser informado no header no seguinte formato: `API_KEY: <TOKEN>` 
> As configurações de limite do token de acesso devem se sobrepor as do IP. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.

#### 🗂️ Estrutura do Projeto
    .
    ├── cmd                  # Entrypoints da aplicação
    │    └── app_rl
    │           └── main.go  ### Entrypoint da aplicação exemplo (que faz uso do rate limiter)
    ├── config               # helpers para configuração da aplicação (viper)
    ├── internal
    │    └── infra           # Implementações de repositórios e conexões com serviços externos
    │           └── web      ### Funções utilitárias
    │                ├── handler           ### Handlers utilizados pelos endpoints da aplicação exemplo
    │                └── webserver.go      ### Implmentação do servidor web
    ├── pkg                  # Pacotes reutilizáveis utilizados na aplicação
    │    └── fcrl            # Rate Limiter
    │           ├── cache       ### Implementações de cache: Redis e Memória
    │           ├── helpers     ### Funções utilitárias
    │           ├── middleware  ### Implementações do middleware de rate limit
    │           ├── rllog       ### Implementações de log para o rate limit
    │           └── limiter.go  ### Implementação de lógica do limiter
    ├── test                 # Testes automatizados
    ├── Dockerfile           # Arquivo de configuração do Docker da aplicação exemplo
    ├── .env                 # Arquivo de parametrizações globais
    └── README.md

#### 🧭 Parametrização

```dotenv
# .env

##> Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PWD=password!
##< Redis

##> Rate Limit
# Rate Limit Type: Second | Minute | Hour | Day
RATE_LIMIT_BY=Minute

# Requests Limit
IP_LIMIT_RATE=4
API_TOKEN_LIMIT_RATE=5

# Time in Second | Minute | Hour | Day according to RATE_LIMIT_BY
IP_WINDOW_TIME=8
API_TOKEN_WINDOW_TIME=10
##< Rate Limit

##> App
APP_PORT=8080
##< App
```

#### 🚀 Execução:
Para executar a aplicação em ambiente local, basta utilizar o docker-compose disponível na raiz do projeto. Para isso, execute o comando abaixo:
```bash
$ docker-compose up # -d (para executar em background)
```

> 💡 **Portas necessárias:**
> - Aplicação: 8080
> - Redis: 6379

#### 📝 Endpoints para validação:

```shell
# Requisição 1 - Sem Token

curl --location 'http://localhost:8080/time/japanese-greetings' \
--header 'Content-Type: application/json'
```

```shell
# Requisição 2 - Com Token

curl --location 'http://localhost:8080/time/greetings' \
--header 'Content-Type: application/json' \
--header 'API_KEY: xQxGOkwWAxhlaaHJorRnhMwpZ1q8xAlIFHnOkDoEqtNiB5NdYXINVthDbmQIyAE2UIIbYt1SqQ4Ych0MG9EJLftbulvm6gH5IqO9e18MWTyAwcmQkBD4ZW8OMvRgKxnu'
```

#### ✍️ Exemplo de utilização do Rate Limiter

```gotemplate
// new server
r := chi.NewRouter()

r.Use(middleware.RateLimiter(
    RateLimitBy,
    middleware.WithIPRateLimiter(IpLimitRate, IpWindowTime),
    middleware.WithApiKeyRateLimiter(ApiTokenLimitRate, ApiTokenWindowTime),
	// Opções de cache (apenas 1 é necessária):
    middleware.WithRedisCache(rediHost, redisPort, redisPwd, ctx),
    middleware.WithRedisClient(redisClient, ctx), // *redis.Client
    middleware.WithMemoryCache(ctx),
))
```

#### 🧪 Teste:

Para executar o teste, basta executar o comando abaixo:

```bash
$ go test -v github.com/tiagoncardoso/fc-pge-rate-limit/test/e2e
```