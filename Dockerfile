FROM golang:1.25.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

COPY .env .env

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

#final image creation
FROM scratch

COPY --from=builder /app/main /

EXPOSE 8000

# Set entrypoint to the binary
ENTRYPOINT ["/main"]
