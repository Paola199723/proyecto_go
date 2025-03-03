compilar en develop
export GO_ENV=develop && go run cmd/main.go

compilar en produccion 
export GO_ENV=develop && go run cmd/main.go

ejecutar Docker
docker build -t miapp .