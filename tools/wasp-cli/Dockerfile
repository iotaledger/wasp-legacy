# syntax=docker/dockerfile:1
ARG GOLANG_IMAGE_TAG=1.24-bullseye

# Build stage
FROM golang:${GOLANG_IMAGE_TAG} AS build
ARG BUILD_LD_FLAGS="--X=github.com/iotaledger/wasp/components/app.Version=v0.0.0-testing"

LABEL org.label-schema.description="Sandbox Wasp CLI"
LABEL org.label-schema.name="iotaledger/sandbox-wasp-cli"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.vcs-url="https://github.com/iotaledger/wasp"

# Ensure ca-certificates are up to date
RUN update-ca-certificates

# Set the current Working Directory inside the container
RUN mkdir /scratch
WORKDIR /scratch

# Prepare the folder where we are putting all the files
RUN mkdir /app

# Make sure that modules only get pulled when the module file has changed
COPY go.mod go.sum ./

RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/root/go/pkg/mod \
  go mod download

# Project build stage
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/root/go/pkg/mod \
  cd ./tools/wasp-cli && go build -o /app/wasp-cli -a -ldflags=${BUILD_LD_FLAGS} .

############################
# Image
############################
# https://console.cloud.google.com/gcr/images/distroless/global/cc-debian11
# using distroless cc "nonroot" image, which includes everything in the base image (glibc, libssl and openssl)
FROM gcr.io/distroless/cc-debian11:nonroot

# Copy the app dir into distroless image
COPY --chown=nonroot:nonroot --from=build /app /app

WORKDIR /app
USER nonroot

ENTRYPOINT ["/app/wasp-cli"]
