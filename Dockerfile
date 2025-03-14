FROM --platform=$BUILDPLATFORM golang:1.24.1-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src/app

COPY go* .
RUN go mod download

COPY . .
RUN go vet -v github.com/wavy-cat/petpet-go/cmd/app
RUN go test -v github.com/wavy-cat/petpet-go/cmd/app

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o server github.com/wavy-cat/petpet-go/cmd/app

FROM gcr.io/distroless/static-debian12
LABEL authors="wavycat"

WORKDIR /app
COPY --from=builder /src/app /app

# Only for Docker Desktop
EXPOSE 3000

ENTRYPOINT ["./server"]
