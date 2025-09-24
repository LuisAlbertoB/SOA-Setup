# --- Etapa de Construcción (Builder) ---
# Usamos una imagen oficial de Go como base para compilar nuestro código.
# Especificamos la versión para asegurar compilaciones consistentes.
FROM golang:1.21-alpine AS builder

# Establecemos el directorio de trabajo dentro del contenedor.
WORKDIR /app

# Copiamos los archivos de dependencias primero.
# Esto aprovecha el cache de Docker: si go.mod/go.sum no cambian,
# no se volverán a descargar las dependencias en futuras construcciones.
COPY go.mod go.sum ./
RUN go mod download

# Copiamos todo el código fuente de la aplicación.
COPY . .

# Compilamos la aplicación.
# -o /app/main: Especifica que el binario ejecutable se llamará 'main' y se guardará en /app.
# CGO_ENABLED=0: Deshabilita CGO para crear un binario estático, lo que lo hace más portable
# y no depende de librerías C del sistema operativo base.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./main.go

# --- Etapa Final (Runner) ---
# Usamos una imagen base mínima (Alpine Linux) para la imagen final.
# Esto reduce drásticamente el tamaño de la imagen y la superficie de ataque.
FROM alpine:latest

# Copiamos el binario compilado desde la etapa 'builder'.
COPY --from=builder /app/main /app/main

# Exponemos el puerto que la aplicación usará (según tu config).
EXPOSE 8080

# El comando que se ejecutará cuando el contenedor inicie.
CMD ["/app/main"]