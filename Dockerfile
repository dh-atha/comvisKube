FROM golang:1.23-alpine AS build

WORKDIR /src

COPY go.mod ./
COPY main.go ./
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/cloudnative-app ./main.go

FROM alpine:3.20

RUN adduser -D -g '' appuser && mkdir -p /app /data/metadata && chown -R appuser:appuser /app /data

USER appuser

WORKDIR /app

COPY --from=build --chown=appuser:appuser /out/cloudnative-app /app/cloudnative-app

EXPOSE 8080

ENTRYPOINT ["/app/cloudnative-app"]