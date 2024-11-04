# Geo IP Info

Geo IP Info (GIP) es un CLI & API que permite obtener información de una dirección IP.

## Prerrequisitos

- [Go 1.23 o superior](https://go.dev/dl/)
- [Docker](https://docs.docker.com/desktop/)
- Docker Compose
- Make (Opcional)

## Instalación

### Opción 1: Instalación con Docker (API)

1. Clonar el repositorio
```bash
git clone https://github.com/cgiraldoz/geo-ip-info.git
````

2. Crear el archivo `.env` en la raíz del proyecto con el siguiente contenido:

```env
FIXER_API_KEY=api_key
IPAPI_API_KEY=api_key
```

3. Ejecutar el comando `docker-compose up -d`
4. Acceder a la URL `http://localhost:3000/api/ip/8.8.8.8` en tu navegador o cliente HTTP
5. ¡Listo!

### Opción 2: Instalación con binarios (API/CLI)

1. Navegar a la [página de release](https://github.com/cgiraldoz/geo-ip-info/releases/tag/v1.0.0)
2. Descargar el archivo correspondiente a tu sistema operativo
3. Descomprimir el archivo
4. Crear el archivo `.env` en la raíz de la carpeta con el siguiente contenido:

```env
FIXER_API_KEY=api_key
IPAPI_API_KEY=api_key
```

5. Ejecutar el comando `docker-compose up -d` en la raíz de la carpeta
6. Para ejecutar el CLI, escribe el comando `./gip ip 8.8.8.8`
7. Para ejecutar la API, escribe el comando `./gip api`
8. Si ejecutas la API, accede a la URL `http://localhost:3000/api/ip/8.8.8.8` en tu navegador o cliente HTTP.
9. ¡Listo!

## Uso

### CLI
Obtener información de una dirección IP
```bash
./gip ip [mi_ip]
```

Ejemplo:
```bash
./gip ip 8.8.8.8
```

Consultar estadísticas de uso
```bash
./gip stats
```

Inicializar la API
```bash
./gip api
```

### API

Ejecutar el comando `./gip api` para inicializar la API. Accede a la URL `http://localhost:3000`


#### Endpoints

- `/api/ip/{ip}`: Obtiene información de una dirección IP
- `/api/stats`: Obtiene las estadísticas de uso

Ejemplo:
```bash
http://localhost:3000/api/ip/8.8.8.8
```

## URLs de interés

- [Fixer - Foreign exchange rates and currency conversion JSON API](https://fixer.io/)
- [IPAPI - IP Address Location API](https://ipapi.com/)
