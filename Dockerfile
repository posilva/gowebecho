FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY * ./

RUN go build -o /service

EXPOSE 4001

CMD ["/service"]
