# Build stage
FROM golang:latest AS builder

WORKDIR /build

# Copiar go.mod e go.sum para aproveitar o cache
COPY ../src/go.mod ../src/go.sum ./

# Copiar dependências locais
COPY ../src/controllers ./controllers
COPY ../src/library ./library
COPY ../src/metrics ./metrics
COPY ../src/models ./models
COPY ../src/whatsapp ./whatsapp
COPY ../src/whatsmeow ./whatsmeow

# Copiar o script de entrypoint
COPY ../docker/docker-entrypoint.sh ./docker-entrypoint.sh

# Baixar as dependências
RUN go mod download

# Copiar o restante dos arquivos do projeto
COPY ../src/ ./src/

# Construir o binário
RUN cd src && go build -o /build/service main.go

# Run stage (imagem final)
FROM golang:latest AS final

WORKDIR /opt/quepasa/

# Copiar o binário e o entrypoint
COPY --from=builder /build/service /opt/quepasa/
COPY --from=builder /build/docker-entrypoint.sh /opt/quepasa/

# Ajustar permissões
RUN chmod +x /opt/quepasa/docker-entrypoint.sh

# Definir entrypoint e comando padrão
ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["sh"]
