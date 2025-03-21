### Build app in docker
FROM golang:alpine AS builder

WORKDIR /src
COPY ./views /src/views
COPY ./models /src/models
COPY ./go.mod /src/go.mod
COPY ./go.sum /src/go.sum
COPY ./main.go /src/main.go

RUN go mod download && go mod verify
###RUN apk update && apk add librdkafka-dev pkgconf
#RUN go install github.com/swaggo/swag/cmd/swag@latest
#RUN swag init
RUN go build -o /out/app .

### Build image
FROM alpine:3

RUN apk update && apk add --no-cache --update bash openssl ca-certificates curl
#RUN apk update && apk add --no-cache --update bash openssl ca-certificates librdkafka-dev pkgconf

WORKDIR /app

COPY --from=builder /out/app .
#COPY --from=builder /src/.env .

### Unused view at this moment
COPY --from=builder /src/views/ ./views

EXPOSE 3000
ENTRYPOINT ["/app/app"]