FROM golang:1.21.5

WORKDIR /app

COPY . .
RUN go mod download

WORKDIR /app/cmd
RUN GOOS=linux go build -o wallet-api-service

EXPOSE 8080

CMD ["./wallet-api-service"]
