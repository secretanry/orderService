FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/go-binary .

FROM alpine:3.21

WORKDIR /app
RUN apk update && apk upgrade && apk add --no-cache
COPY --from=builder /app/go-binary /app/go-binary
COPY --from=builder /app/templates ./templates
COPY docs/swagger.json docs/swagger.json
CMD ["/app/go-binary"]