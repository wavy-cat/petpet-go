FROM --platform=$BUILDPLATFORM golang:1.24.1-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src/app

COPY go* .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o server github.com/wavy-cat/petpet-go/cmd/app

FROM gcr.io/distroless/static-debian12:nonroot

LABEL authors="WavyCat"
LABEL org.opencontainers.image.source="https://github.com/wavy-cat/petpet-go"

WORKDIR /app
COPY --from=builder /src/app /app

USER 1002

# Only for Docker Desktop
EXPOSE 3000

ENTRYPOINT ["./server"]
