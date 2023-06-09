FROM golang:1.20.2-alpine as builder
WORKDIR /src
COPY . .

# This is where one could build the application code as well.
RUN go build ./cmd/chatgw


# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM airdb/base:latest

WORKDIR /app
RUN apk update && apk add ca-certificates iptables ip6tables && rm -rf /var/cache/apk/*

# Copy binary to production image
COPY --from=builder /src/chatgw /app/chatgw

# Run on container startup.
EXPOSE 30120
ENTRYPOINT ["/app/chatgw"]
