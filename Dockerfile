FROM golang:latest

WORKDIR /var/www/app

COPY . .

USER taylor

CMD ["run", "-cors-trusted-origins='http://localhost:3000 http://localhost:3000/*'", "cmd/api"]

EXPOSE 4000