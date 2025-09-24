#!/bin/bash

# setup.sh: Script para configurar un servidor de despliegue desde cero.
# Instala Docker, Docker Compose y AWS CLI, y genera el archivo docker-compose.yml.

# Detener la ejecución inmediatamente si un comando falla.
set -e

echo "--- Iniciando script de configuración y despliegue para el servidor ---"

# --- Variables de Configuración ---
# Estas variables deben ser configuradas como variables de entorno en el servidor.
# Se leen aquí con valores por defecto para facilitar pruebas locales.
export APP_NAME=${APP_NAME:-rob-api-go}
export SERVER_PORT=${SERVER_PORT:-8080}
export MARIADB_DATABASE=${MARIADB_DATABASE:-appstore_db}
export MARIADB_USER=${MARIADB_USER:-user}
export MARIADB_PASSWORD=${MARIADB_PASSWORD:-password}
export MARIADB_ROOT_PASSWORD=${MARIADB_ROOT_PASSWORD:-rootpassword}
export AWS_REGION=${AWS_REGION:-us-east-1}
export AWS_S3_BUCKET_NAME=${AWS_S3_BUCKET_NAME:-mi-bucket-s3-unico}

# Construir el DSN (Data Source Name) de MariaDB a partir de las variables de entorno.
export MARIADB_DSN="${MARIADB_USER}:${MARIADB_PASSWORD}@tcp(mariadb:3306)/${MARIADB_DATABASE}?charset=utf8mb4&parseTime=True&loc=Local"

# --- 1. Instalación de Prerrequisitos ---

install_docker() {
    if ! command -v docker &> /dev/null; then
        echo "Docker no encontrado. Instalando..."
        apt-get update
        apt-get install -y apt-transport-https ca-certificates curl software-properties-common
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
        echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
        apt-get update
        apt-get install -y docker-ce docker-ce-cli containerd.io
        # Añadir el usuario 'ubuntu' (común en EC2) al grupo de docker para evitar usar sudo con docker
        usermod -aG docker ubuntu || true 
        echo "Docker instalado correctamente."
    else
        echo "Docker ya está instalado. Saltando instalación."
    fi
}

install_compose() {
    # La versión 1.29.2 es una versión estable y conocida. Se puede actualizar si es necesario.
    COMPOSE_VERSION="1.29.2"
    if ! command -v docker-compose &> /dev/null; then
        echo "Docker Compose no encontrado. Instalando..."
        curl -L "https://github.com/docker/compose/releases/download/${COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        chmod +x /usr/local/bin/docker-compose
        echo "Docker Compose instalado correctamente."
    else
        echo "Docker Compose ya está instalado. Saltando instalación."
    fi
}

install_aws_cli() {
    if ! command -v aws &> /dev/null; then
        echo "AWS CLI no encontrado. Instalando..."
        apt-get install -y unzip
        curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
        unzip awscliv2.zip
        ./aws/install
        rm -rf aws awscliv2.zip
        echo "AWS CLI instalado correctamente."
    else
        echo "AWS CLI ya está instalado. Saltando instalación."
    fi
}

# --- 2. Creación del archivo Docker Compose ---
create_docker_compose() {
    echo "--- Creando/actualizando archivo docker-compose.yml ---"
    # Usamos un "here document" para escribir el contenido del archivo.
    # Las variables de entorno ($VAR) serán sustituidas por sus valores actuales.
    # El archivo se creará en el directorio home del usuario que ejecute el script.
    cat > ./docker-compose.yml << EOL
version: '3.8'

services:
  mariadb:
    image: mariadb:10.6
    container_name: mariadb_db
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: \${MARIADB_ROOT_PASSWORD}
      MYSQL_DATABASE: \${MARIADB_DATABASE}
      MYSQL_USER: \${MARIADB_USER}
      MYSQL_PASSWORD: \${MARIADB_PASSWORD}
    volumes:
      - mariadb_data:/var/lib/mysql
    networks:
      - app-network

  app:
    # La imagen se especificará dinámicamente en el comando de despliegue
    # image: ACCOUNT_ID.dkr.ecr.REGION.amazonaws.com/rob-api-go:TAG
    container_name: \${APP_NAME}
    restart: always
    ports:
      - "\${SERVER_PORT}:\${SERVER_PORT}"
    environment:
      - SERVER_PORT=\${SERVER_PORT}
      - MARIADB_DSN=\${MARIADB_DSN}
      - AWS_REGION=\${AWS_REGION}
      - AWS_S3_BUCKET_NAME=\${AWS_S3_BUCKET_NAME}
      # Esto es crucial para que el SDK de Go dentro del contenedor
      # pueda asumir el rol IAM de la instancia EC2.
      - AWS_SDK_LOAD_CONFIG=1
    depends_on:
      - mariadb
    networks:
      - app-network

volumes:
  mariadb_data:

networks:
  app-network:
EOL
    echo "docker-compose.yml creado/actualizado en el directorio actual."
}

# --- Ejecución Principal ---

# El script debe ejecutarse con privilegios de superusuario (sudo) para instalar software.
if [ "$EUID" -ne 0 ]; then 
  echo "Por favor, ejecuta este script como root o con sudo."
  exit 1
fi

echo "--- Verificando e instalando prerrequisitos ---"
install_docker
install_compose
install_aws_cli

create_docker_compose

echo "--- Configuración inicial del servidor completada ---"
echo "El servidor está listo para recibir despliegues desde el pipeline de CI/CD."