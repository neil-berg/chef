FROM golang:alpine as builder

ENV GO111MODILE=on

RUN apk update && apk add --no-cache git

WORKDIR /app

## Install dependencies, copy source code, then build binary
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

## Slimmer production image to hold the built binary
FROM alpine:latest
WORKDIR /usr/src/app
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]