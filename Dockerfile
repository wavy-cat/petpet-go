FROM --platform=$BUILDPLATFORM golang:1.25.6-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src/app

COPY go* .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags="-s -w" -o petpet-go github.com/wavy-cat/petpet-go/cmd/app

FROM gcr.io/distroless/static-debian13

LABEL authors="WavyCat"
LABEL org.opencontainers.image.source="https://github.com/wavy-cat/petpet-go"

WORKDIR /app
COPY --from=builder /src/app /app

ENTRYPOINT ["./petpet-go"]
