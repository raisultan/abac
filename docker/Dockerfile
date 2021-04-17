FROM golang:1.16.3-alpine

WORKDIR /app

COPY ./app/go.mod /app
COPY ./app/go.sum /app

RUN apk add --update make \
    && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
    && go mod download

COPY ./app /app

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
