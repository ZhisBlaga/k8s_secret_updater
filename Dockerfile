FROM --platform=${BUILDPLATFORM} docker.io/golang:1.22.0-alpine@sha256:8e96e6cff6a388c2f70f5f662b64120941fcd7d4b89d62fec87520323a316bd9 AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /main
FROM alpine:3
COPY --from=builder main /bin/main
ENTRYPOINT ["/bin/main"]