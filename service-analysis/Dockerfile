# --- Estágio 1: Build ---

FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /server ./cmd/analysis/main.go


# --- Estágio 2: Final ---

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /server .


# Porta que o serviço usa (interno docker)
EXPOSE 8082

CMD ["./server"]