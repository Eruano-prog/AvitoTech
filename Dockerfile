FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main /app/cmd/server/

CMD ["./main"]