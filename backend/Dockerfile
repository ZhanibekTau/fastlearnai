FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o fastlearn-api ./cmd/main.go

FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/fastlearn-api /usr/local/bin/
EXPOSE 8080
CMD ["fastlearn-api"]
