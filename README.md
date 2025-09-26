# Automatización de CI/CD con Docker y AWS

Este proyecto es un ejemplo de cómo configurar un flujo de trabajo completo de Integración y Despliegue Continuo (CI/CD) para una aplicación Go. El objetivo es automatizar todo el proceso desde que un desarrollador integra su código hasta que la nueva versión está funcionando en un servidor en la nube.

## ¿Cómo funciona el Pipeline de CI/CD?

El pipeline se encarga de todo el trabajo pesado de forma automática. Lo hemos configurado para que se active solo cuando un **Pull Request es aprobado y fusionado** en la rama `develop`. Así nos aseguramos de que solo el código revisado llegue al servidor.

El proceso se divide en dos grandes pasos:

1.  **Construir y Publicar (`build-and-push`)**:
    *   **Activación**: Un desarrollador fusiona su código en `develop`.
    *   **Construcción**: GitHub Actions toma el código y usa el `Dockerfile` para construir una imagen Docker de la aplicación. Esta imagen es como una "caja" sellada que contiene todo lo que la app necesita para funcionar.
    *   **Publicación**: Una vez creada, la imagen se sube a nuestro registro de contenedores privado en **Amazon ECR**. Cada imagen lleva una etiqueta única (el ID del commit) para saber exactamente qué versión del código contiene.

2.  **Desplegar en el Servidor (`deploy-to-ec2`)**:
    *   **Conexión**: Una vez que la imagen está guardada en ECR, el pipeline se conecta de forma segura a nuestro servidor **EC2** en AWS usando SSH.
    *   **Actualización**: Ya en el servidor, ejecuta una serie de comandos:
        *   Se autentica en Amazon ECR para poder descargar la imagen.
        *   Usa `docker-compose` para descargar la nueva versión de la imagen que acabamos de subir.
        *   Reinicia el contenedor de la aplicación para que empiece a usar el nuevo código. ¡Todo esto sin afectar al contenedor de la base de datos!

En resumen: **código fusionado -> imagen creada -> imagen publicada -> servidor actualizado**. Todo en cuestión de minutos y sin intervención manual.

## ¿Y las Migraciones de la Base de Datos?

¡Buena pregunta! No queremos actualizar la base de datos a mano cada vez que desplegamos.

La solución está dentro de la propia aplicación. Usamos una librería llamada **GORM** que tiene una función muy útil: `AutoMigrate`.

Así es como funciona:

1.  Cuando el pipeline despliega el nuevo contenedor y la aplicación se inicia, una de las primeras cosas que hace el código es conectarse a la base de datos.
2.  Justo después de conectar, se ejecuta `DB.AutoMigrate(...)`.
3.  Esta función compara los modelos de datos que tenemos definidos en el código Go (como las estructuras `User`, `App`, etc.) con las tablas que existen en la base de datos MariaDB.
4.  Si `AutoMigrate` detecta que falta una tabla o una columna nueva, **la crea automáticamente**.

De esta forma, la base de datos siempre está sincronizada con la versión del código que se está ejecutando, y las migraciones se aplican solas en cada despliegue.

## Evidencias del Despliegue

A continuación, se muestran algunas capturas que demuestran el pipeline en acción:

**1. Pipeline de GitHub Actions completado con éxito:**
!Pipeline Exitoso

**2. Nueva imagen publicada en Amazon ECR:**
!Imagen en ECR

**3. Logs del despliegue en la instancia EC2:**
!Logs de Despliegue

## Conclusiones Personales

*(Aquí puedes añadir tus conclusiones sobre el proceso, los desafíos encontrados y las ventajas de usar este enfoque de CI/CD con contenedores y la nube).*
