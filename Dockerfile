ARG GOLANG_IMAGE_TAG=1.18.2-alpine3.16
ARG ALPINE_IMAGE_TAG=3.16.0

FROM golang:$GOLANG_IMAGE_TAG AS base

WORKDIR /app

COPY ./gruppetto/go.mod ./
COPY ./gruppetto/go.sum ./
RUN go mod download

COPY ./gruppetto/ ./

FROM base as build

RUN go build -o ./gruppetto.go ./cmd/main.go

FROM base as test

ENV CGO_ENABLED=0

ENTRYPOINT ["go",  "test", "./..."]

FROM alpine:$ALPINE_IMAGE_TAG as production

WORKDIR /app

COPY --from=build /app/gruppetto.go ./gruppetto.go

ENV GRUPPETTO_IP=127.0.0.1 \
    GRUPPETTO_PORT=3478 \
    GRUPPETTO_USER=user \
    GRUPPETTO_PASSWORD=password \
    GRUPPETTO_REALM=gruppetto

EXPOSE 8080

ENTRYPOINT ["./gruppetto.go"]
