FROM node:latest as builder

WORKDIR /swag
COPY . .
RUN npm install -g multi-file-swagger
RUN multi-file-swagger -o yaml index.yaml > pasmas.yaml


FROM swaggerapi/swagger-ui

COPY --from=builder /swag/pasmas.yaml /swag/pasmas.yaml

ENV SWAGGER_JSON=/swag/pasmas.yaml

EXPOSE 8000
