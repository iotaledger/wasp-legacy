# syntax=docker/dockerfile:1
ARG GOLANG_IMAGE_TAG=1.24-bullseye

# Build stage
FROM golang:${GOLANG_IMAGE_TAG} AS build
ARG BUILD_TAGS=rocksdb
ARG BUILD_LD_FLAGS="--X=github.com/iotaledger/wasp/components/app.Version=v0.0.0-testing"

LABEL org.label-schema.description="Wasp"
LABEL org.label-schema.name="iotaledger/wasp"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.vcs-url="https://github.com/iotaledger/wasp"

# Ensure ca-certificates are up to date
RUN update-ca-certificates

# Set the current Working Directory inside the container
RUN mkdir /scratch
WORKDIR /scratch

# Prepare the folder where we are putting all the files
RUN mkdir /app
RUN mkdir /app/waspdb

# Make sure that modules only get pulled when the module file has changed
COPY go.mod go.sum ./

RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/root/go/pkg/mod \
  go mod download

# Project build stage
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/root/go/pkg/mod \
  go build -o /app/wasp -a -tags=${BUILD_TAGS} -ldflags=${BUILD_LD_FLAGS} .

############################
# Image
############################
# https://console.cloud.google.com/gcr/images/distroless/global/cc-debian11
# using distroless cc "nonroot" image, which includes everything in the base image (glibc, libssl and openssl)
FROM gcr.io/distroless/cc-debian11:nonroot

EXPOSE 9090/tcp
EXPOSE 6060/tcp
EXPOSE 4000/udp
EXPOSE 4000/tcp

HEALTHCHECK --interval=10s --timeout=5s --retries=30 CMD ["/app/wasp", "tools", "node-health"]

# Copy the app dir into distroless image
COPY --chown=nonroot:nonroot --from=build /app /app
COPY --chown=nonroot:nonroot --from=build /app/waspdb /app/waspdb

WORKDIR /app
USER nonroot

ENTRYPOINT ["/app/wasp"]
