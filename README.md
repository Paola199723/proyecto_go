# Proyecto: Sistema de Autenticación y Recomendaciones en Go

## Descripción
Este proyecto es una API REST desarrollada en Go utilizando el framework Gin. Permite la autenticación de usuarios, la obtención de listas de acciones y la generación de recomendaciones financieras.

## Tecnologías utilizadas
- **Go** (Golang)
- **Gin** (Framework web para Go)
- **GORM** (ORM para la gestión de la base de datos)
- **PostgreSQL** (Base de datos relacional)

## Instalación
1. Clona este repositorio:
   ```sh
   git clone https://github.com/tu_usuario/proyecto_go.git
   cd proyecto_go
   ```

2. Instala las dependencias:
   ```sh
   go mod tidy
   ```

3. Configura la base de datos en `configuration/` y establece las credenciales correctas.

## Uso
### Ejecutar la API
Para iniciar el servidor, ejecuta el siguiente comando:
```sh
 go run main.go
```
Por defecto, el servidor se ejecutará en `http://localhost:8080`.

### Endpoints disponibles
#### Autenticación de usuario
- **POST** `/user/login`
  - Entrada: JSON con `username` y `password`.
  - Salida: Token de autenticación.

#### Listado de acciones
- **GET** `/user/stocks`
  - Requiere el token de autenticación en el header `Authorization`.
  - Retorna una lista de acciones con su información.

#### Recomendaciones
- **GET** `/user/recomendations`
  - Retorna las tres mejores recomendaciones de inversión.


## Autor
Desarrollado por Paola casadiegos Vaca.
