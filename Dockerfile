FROM golang:1.22.5 AS builder

WORKDIR /api

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin .

FROM alpine:3.20

WORKDIR /api

COPY --from=builder /api/bin /api/bin

RUN apk --no-cache add ca-certificates

ENTRYPOINT [ "/api/bin" ]