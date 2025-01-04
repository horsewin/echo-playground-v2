# Multi stage building strategy for reducing image size.
FROM golang:1.23.4 AS build-env
ENV GO111MODULE=on \
    GOPATH=/go \
    GOBIN=/go/bin \
    PATH=/go/bin:$PATH

# Set working directory
RUN mkdir /app
WORKDIR /app

# Install each dependencies
COPY go.mod go.sum ./
RUN go mod download

# Install golangci-lint
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4

# COPY main module
COPY . /app

# Check and Build
RUN make validate && \
    make build-linux

### If use TLS connection in container, add ca-certificates following command.
### > RUN apt-get update && apt-get install -y ca-certificates
FROM gcr.io/distroless/base-debian10
COPY --from=build-env /app/bin/main /
EXPOSE 80
ENTRYPOINT ["/main"]
