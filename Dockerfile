###############################################################################
# Temporary Builder container
###############################################################################
# Official golang:1.15-alpine3.12 image from Docker hub.
# Alpine FROM golang@sha256:59eae48746048266891b7839f7bb9ac54a05cec6170f17ed9f4fd60b331b644b as builder
FROM golang@sha256:5219b39d2d6bf723fb0221633a0ff831b0f89a94beb5a8003c7ff18003f48ead as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates are required to call HTTPS endpoints.
# Alpine RUN apk update && apk add --no-cache git ca-certificates docker ps tzdata && update-ca-certificates
RUN apt-get update && apt-get install -y git ca-certificates docker tzdata && update-ca-certificates

# Create appuser, to avoid running as root inside container.
ENV USER=appuser
ENV UID=10001

# Do not assign a password
# GECOS field
# Set homedir, but do not create it
# Set (no) shell
# Set user ID
# Alpine RUN adduser \
# Alpine     -D \
# Alpine     -g "" \
# Alpine     -H -h "/nonexistent" \
# Alpine     -s "/sbin/nologin" \
# Alpine     -u "${UID}" \
# Alpine     "${USER}"
RUN adduser \
    --disabled-password \
    --gecos "" \
    --no-create-home --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --uid "${UID}" \
    "${USER}"

WORKDIR $GOPATH/src/github.com/otaviokr/mongodb-proxy-ms
COPY . .

# Fetch dependencies.
RUN go get -d -v

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o /go/bin/mongodb-proxy-ms .

###############################################################################
# Build the final, smaller image
###############################################################################
FROM scratch

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable.
COPY --from=builder /go/bin/mongodb-proxy-ms /go/bin/mongodb-proxy-ms

# Switch to the non-root user we created.
USER appuser:appuser

# Required information to connect to database
ENV GIN_MODE=${GIN_MODE:-"release"}
ENV	MONGODB_HOST=${MONGODB_HOST:-"localhost"}
ENV MONGODB_PORT=${MONGODB_PORT:-"27017"}
ENV MONGODB_USER=${MONGODB_USER}
ENV MONGODB_PASS=${MONGODB_PASS}

EXPOSE 8080

# Run our executable.
ENTRYPOINT ["/go/bin/mongodb-proxy-ms"]
