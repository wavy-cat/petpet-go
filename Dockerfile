FROM golang:1.23.3-alpine AS builder

WORKDIR /src/app

COPY . .

RUN go mod download
RUN go build -o petpet github.com/wavy-cat/petpet-go/cmd/app

FROM alpine
LABEL authors="wavycat"

WORKDIR /app
COPY --from=builder /src/app /app

EXPOSE 80

ENTRYPOINT ["./petpet"]