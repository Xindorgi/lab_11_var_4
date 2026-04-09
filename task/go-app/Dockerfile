FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod .
COPY main.go .
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o /app .

FROM scratch
COPY --from=builder /app /app
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app", "-health"]
ENTRYPOINT ["/app"]