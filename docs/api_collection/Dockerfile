FROM node:18 as compiler

WORKDIR /app/api
COPY ./Api.json .
COPY ./api_options.json .

RUN npm i postman-to-openapi -g
RUN p2o ./Api.json -f ./index.yaml -o ./api_options.json


FROM redocly/cli as builder

WORKDIR /usr/share/openapi
COPY --from=compiler /app/api/index.yaml .

RUN redocly build-docs index.yaml --output ./index.html 

FROM redocly/redoc

COPY --from=builder /usr/share/openapi/index.html /usr/share/nginx/html/index.html
ENV SPEC_URL=/usr/share/nginx/html/index.html

EXPOSE 80

