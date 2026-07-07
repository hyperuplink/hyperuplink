# syntax=docker/dockerfile:1

# ---[ BUILD ]---------------------------------------------------------------- #
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS="$TARGETOS" \
    GOARCH="$TARGETARCH" \
    GOARM="${TARGETVARIANT#v}" \
    go build -trimpath -buildvcs=false \
      -ldflags="-s -w \
        -X xn--gckvb8fzb.com/hyperuplink/runtime.Version=${VERSION} \
        -X xn--gckvb8fzb.com/hyperuplink/runtime.Commit=${COMMIT} \
        -X xn--gckvb8fzb.com/hyperuplink/runtime.Date=${DATE}" \
      -o /out/hyperuplink .

# ---[ RUNTIME ]-------------------------------------------------------------- #
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -g 10001 hyperuplink \
    && adduser -D -H -u 10001 -G hyperuplink hyperuplink \
    && mkdir -p /var/lib/hyperuplink/media \
    && chown -R 10001:10001 /var/lib/hyperuplink

COPY --from=builder /out/hyperuplink /usr/local/bin/hyperuplink
COPY hyperuplink.toml /etc/hyperuplink.toml

# Run rootless (numeric UID:GID works with Podman and arbitrary-UID remapping)
USER 10001:10001

# Non-privileged port so rootless Podman can bind it without extra capabilities
EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/hyperuplink"]
CMD ["-c", "file:///etc/hyperuplink.toml"]
