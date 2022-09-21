FROM golang:latest

WORKDIR /var/www/app

COPY . .

RUN go build ./cmd/api

CMD ["./api", "-cors-trusted-origins='http://localhost:3000 http://localhost:3000/*'"]

EXPOSE 4000