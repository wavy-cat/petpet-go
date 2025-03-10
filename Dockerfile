FROM golang:1.24.1-alpine AS builder

WORKDIR /src/app

COPY go* .
RUN go mod download

COPY . .
RUN go vet -v github.com/wavy-cat/petpet-go/cmd/app
RUN go test -v github.com/wavy-cat/petpet-go/cmd/app

RUN CGO_ENABLED=0 go build -o petpet github.com/wavy-cat/petpet-go/cmd/app

FROM gcr.io/distroless/static-debian12
LABEL authors="wavycat"

WORKDIR /app
COPY --from=builder /src/app /app

# Only for Docker Desktop
EXPOSE 3000

ENTRYPOINT ["./petpet"]
