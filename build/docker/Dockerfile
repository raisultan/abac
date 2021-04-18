FROM golang:1.16.3-alpine

WORKDIR /abac

COPY ./go.mod /abac
COPY ./go.sum /abac

RUN apk add --update make \
    && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
    && go mod download

COPY ./ /abac

RUN go build -o main ./cmd/server/

EXPOSE 8080

CMD ["./main"]
