# 1. BUILD EXECUTABLE BINARY
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go build -o order-service

# 2. BUILD A SMALL IMAGE
FROM alpine
WORKDIR /app
COPY --from=builder /app/order-service /app/
CMD [ "./order-service" ]