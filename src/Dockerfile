FROM golang as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o pasmas .

EXPOSE 8080


FROM alpine
WORKDIR /app
COPY --from=builder /app/pasmas /app/pasmas
COPY ./database/sqlScripts/ /app/database/sqlScripts/
RUN chmod +x /app/pasmas

EXPOSE 8080
CMD ["/app/pasmas"]
