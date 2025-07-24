# Dockerfile

# ---- Estágio 1: Build ----
# Usamos uma imagem oficial do Go para compilar nosso código.
# Chamamos este estágio de "builder".
FROM golang:1.24-alpine AS builder
# Define o diretório de trabalho dentro do contêiner.
WORKDIR /app

# Copia os arquivos de gerenciamento de dependências.
COPY go.mod ./
COPY go.sum ./

# Baixa as dependências.
RUN go mod download

# Copia todo o código-fonte da nossa aplicação para o contêiner.
COPY . .

# Compila o código Go.
# CGO_ENABLED=0 desativa o CGO para criar um binário estático.
# -o /app/server cria um arquivo de saída chamado "server" no diretório /app.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

# ---- Estágio 2: Final ----
# Agora usamos uma imagem "vazia" (scratch) que não tem NADA.
# Isso torna nossa imagem final extremamente pequena e segura.
# A imagem `alpine` é uma alternativa caso precise de um shell para depuração.
FROM scratch

# Define o diretório de trabalho.
WORKDIR /app

# Copia APENAS o binário compilado do estágio "builder".
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

# Expõe a porta 8080 para que o Docker saiba que nosso app
# escuta nesta porta.
EXPOSE 8080

# Comando para executar quando o contêiner iniciar.
# Executa nosso binário compilado.
CMD ["/app/server"]
