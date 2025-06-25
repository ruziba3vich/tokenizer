# Use a minimal Go image
FROM golang:1.24.4-alpine
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/main.go
EXPOSE 7777
CMD ["./main"]
