# PasMAs

The **Pas**senger **M**anagement **As**sistant (PasMAs) is an open-source web application designed to simplify the management of sightseeing flights on open days for smaller flying clubs. This assistant enables you to organize various flight slots for each plane efficiently. It supports you by handling flight logic calculations such as determining the maximum take-off mass to prevent overloaded flights.

The application consists of both a backend and a frontend. This repository contains the code for the PasMAS backend application. The code for the frontend can be found at [https://github.com/FGK-PASMAS/FGK_PASMAS_frontend](https://github.com/FGK-PASMAS/FGK_PASMAS_frontend).

## Features

## Getting started
The PasMAs backend application can be deployed via Docker.

### Configuration
Before running the backend application, you need to configure it via environment variables. Simply copy the `.env.example` file to a `.env` file, and set the correct configuration values accordingly.

Important Environment Variables

- DATABASE_HOSTNAME: The hostname for the PostgreSQL database.
- DATABASE_USER: The username for the PostgreSQL database.
- DATABASE_PASSWORD: The password for the PostgreSQL database.
- DATABASE_NAME: The name of the PostgreSQL database.
- JWT_ENCODING: The secret key for JWT encoding.
- JWT_ISSUER: The issuer for JWT tokens.
- ADMIN_PASSWORD: The password for the admin user.

For SSL/TLS configuration:

- TLS_CERT_PATH: The path to the TLS certificate on the host.
- TLS_KEY_PATH: The path to the TLS key on the host.

### Deployment
PasMAs comes with 2 included docker-compose files. One is for using the application over HTTP, and the other one for HTTPS. Make sure you have an up-to-date version of [Docker](https://www.docker.com/) installed on your server.

#### without TLS
Navigate to the project directory and run the following commands to run the backend application as a container:

```
docker compose up -d
```

#### with TLS
Navigate to the project directory and run the following commands to run the backend application as a container:
Before running the following command, ensure that you have configured the TLS_CERT_PATH and TLS_KEY_PATH in your .env file and that these paths correspond to the correct locations on your host machine. Then, navigate to the project directory and run the following command:

```
docker-compose -f docker-compose-ssl.yaml up -d
```

### OpenAPI Documentation

In addition to the backend application, a Redocly (OpenAPI) server runs on port 8000, providing the API specification. This can be useful for understanding the available endpoints and their usage. If not needed, this server can be removed from the Docker Compose configuration.

## Additional Notes
Please note that PasMAs is still under development but is currently not actively maintained by the original authors.
