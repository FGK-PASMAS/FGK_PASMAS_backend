FROM redocly/cli as builder

WORKDIR /usr/share/openapi
COPY . .

RUN redocly build-docs index.yaml --output ./index.html --title PASMAS-API

FROM redocly/redoc

COPY --from=builder /usr/share/openapi/index.html /usr/share/nginx/html/index.html
ENV SPEC_URL=/usr/share/nginx/html/index.html

EXPOSE 80
