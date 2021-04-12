FROM golang:1.16.3-alpine

WORKDIR /app

COPY ./app/go.mod /app

RUN go mod download

COPY ./app /app

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
