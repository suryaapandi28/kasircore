# step 1: build stage
FROM golang:1.22 AS builder

WORKDIR /app
COPY . .

# download dep & build binary
RUN go mod tidy
RUN go build -o myapp main.go

# step 2: run stage
FROM debian:bookworm-slim

WORKDIR /root/
COPY --from=builder /app/myapp .

EXPOSE 8080
CMD ["./myapp"]
