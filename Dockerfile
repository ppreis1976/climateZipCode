# Etapa de construção
FROM golang:1.23 AS builder

# Definir o diretório de trabalho
WORKDIR /app

# Copiar os arquivos de dependências
COPY go.mod go.sum ./

# Baixar as dependências
RUN go mod download

# Copiar o código fonte
COPY . .

# Compilar o aplicativo
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o climatezipcode ./cmd/api

# Etapa de execução
FROM scratch

# Copiar o binário compilado da etapa de construção
COPY --from=builder /app/climatezipcode /climatezipcode

# Definir o comando de entrada
ENTRYPOINT ["/climatezipcode"]