# Etapa de compilación
FROM golang:1.24-alpine AS builder

# Instalar dependencias necesarias
RUN apk add --no-cache curl busybox-extras

# Definir directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias primero (optimiza caché de Docker)
COPY go.mod go.sum ./
RUN go mod tidy && go mod download && go mod verify

# Copiar el código fuente
COPY . .

# Verificar si el directorio cmd/ existe antes de compilar
RUN test -d ./cmd/ || (echo "Directorio cmd/ no encontrado" && exit 1)

# Compilar el binario
RUN go build -o main ./cmd/main.go

# Etapa final (ejecución en imagen mínima)
FROM alpine:latest

# Definir directorio de trabajo
WORKDIR /app

# Copiar solo el binario compilado desde la etapa anterior
COPY --from=builder /app/main .

# Permitir ejecución del binario
RUN chmod +x ./main

# Comando por defecto
CMD ["/app/main"]
