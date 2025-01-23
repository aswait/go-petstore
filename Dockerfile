FROM golang:1.20-alpine as builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main ./cmd/api

FROM alpine:latest

COPY --from=builder /app/main /main

COPY --from=builder /app/static /public

EXPOSE 8080

CMD [ "./main" ]