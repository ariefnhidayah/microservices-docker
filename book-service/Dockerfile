FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go build -o book-service

FROM alpine
WORKDIR /app
COPY --from=builder /app/book-service /app/
CMD [ "./book-service" ]