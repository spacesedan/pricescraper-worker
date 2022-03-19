# Build Stage
FROM golang:1.18-alpine3.15 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o pricescraper-worker cmd/main.go

# Run Stage
FROM alpine:3.15

WORKDIR /app

COPY --from=builder /app/pricescraper-worker .

ENV APP_ENV=production
# Local MongoDB instance
# ENV MONGO_URI=mongodb://root:password@147.182.254.160:27017
# Working MongoDB instance
ENV MONGO_URI=mongodb+srv://dbUser:z4uFbtz2H9xJFerr@cluster0.nmuhd.mongodb.net/test?retryWrites=true&tlsInsecure=true
ENV DB=nft-db
## PORT NEEDS TO BE DEFINED ELSEWHERE

ENTRYPOINT ["/app/pricescraper-worker"]
