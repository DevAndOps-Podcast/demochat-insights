#build stage
FROM golang:1.24.6-alpine3.22 AS builder
RUN apk add --no-cache git upx

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download -x

COPY . .

RUN go build \
  -ldflags="-s -w" \
  -o app -v .

RUN upx app

#final stage
FROM alpine:3.22

RUN apk update
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app .

USER nonroot:nonroot
ENTRYPOINT ["./app"]
