FROM golang:1.17-buster as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -v -o server

# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM mcr.microsoft.com/playwright:focal
# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /app/server

# Run the web service on container startup.
ENTRYPOINT ["/app/server"]

# [END run_helloworld_dockerfile]
# [END cloudrun_helloworld_dockerfile]
