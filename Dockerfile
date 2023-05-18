FROM golang:1.18.9-alpine
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/main.go
CMD ["./main"]