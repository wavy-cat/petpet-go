FROM golang:1.23.3-alpine
LABEL authors="wavycat"

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /petpet github.com/wavy-cat/petpet-go/cmd/app

ENTRYPOINT ["/petpet"]