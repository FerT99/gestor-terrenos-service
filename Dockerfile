# Etapa de construcción (Builder)
FROM golang:1.25-alpine AS builder

# Habilitar CGO puede ser necesario si usas ciertos drivers de db, pero para pgx puro no suele ser necesario.
# Por si acaso, usamos alpine con certificados actualizados.
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copiar go.mod y go.sum primero para aprovechar caché
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Etapa final (Producción)
FROM alpine:latest

# Instalar ca-certificates para permitir peticiones HTTPS y tzdata para zonas horarias
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copiar el ejecutable desde el builder
COPY --from=builder /app/main .

# Exponer el puerto por defecto (Railway lo ignorará y usará la variable PORT, pero es buena práctica)
EXPOSE 8080

# Comando para iniciar la aplicación
CMD ["./main"]
