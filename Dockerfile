# builder
FROM golang:1.25.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/bin/app

# runtime
FROM alpine:3.18
COPY --from=builder /app/bin/app /app/app
WORKDIR /app
EXPOSE 8080
CMD ["/app/app"]
