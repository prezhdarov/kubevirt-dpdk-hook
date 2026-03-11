FROM golang:alpine AS builder

# Add ca-certs
#RUN apk add --update --no-cache ca-certificates

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 

ADD . /build
WORKDIR /build

RUN go mod download && \
    go build -a -ldflags '-extldflags "-static"' -o ovs cmd/sidecar/main.go

FROM scratch

COPY --from=builder /build/ovs /usr/bin/ovs

ENTRYPOINT [ "/usr/bin/ovs" ]