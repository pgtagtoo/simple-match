FROM golang:1.17

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

EXPOSE 8080

CMD ["go", "run", "main.go"]
