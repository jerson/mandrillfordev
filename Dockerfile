# syntax=docker/dockerfile:1.6

FROM --platform=$BUILDPLATFORM golang:1.22 AS build
WORKDIR /app

# Pre-fetch deps separately for better layer caching
COPY go.mod .
RUN go mod download

# Copy the rest of the source
COPY . .

# Build for the requested target platform
ARG TARGETOS TARGETARCH
ENV CGO_ENABLED=0
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/mandrill-dev ./cmd/mandrill-dev

FROM --platform=$TARGETPLATFORM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /out/mandrill-dev /mandrill-dev
ENV PORT=8080
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 CMD ["/mandrill-dev","-healthcheck"]
USER 65532:65532
ENTRYPOINT ["/mandrill-dev"]
