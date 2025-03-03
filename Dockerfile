# Imagen base con Go
FROM golang:1.24-alpine AS builder

# Definir directorio de trabajo
WORKDIR /app

# Copiar solo los archivos de dependencias primero (optimiza caché de Docker)
COPY go.mod go.sum ./

# Descargar dependencias antes de copiar el código
RUN go mod tidy && go mod download && go mod verify
RUN apk add --no-cache curl

# Copiar el resto del código del proyecto
COPY . .

# Verificar si cmd/main.go existe antes de compilar (para debugging)
RUN ls -l ./cmd/

# Construir el binario
RUN go build -o main ./cmd/main.go

# Seg
CMD ["/app/main"]
EXPOSE 8080

