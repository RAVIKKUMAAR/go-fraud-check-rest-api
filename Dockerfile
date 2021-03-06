#FROM golang:latest for local build
FROM icr.io/ibmz/golang:1.15

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main" ]
