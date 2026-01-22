FROM golang:alpine AS builder

# Add ca-certs
#RUN apk add --update --no-cache ca-certificates

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 

ADD . /build
WORKDIR /build

RUN go mod download && \
    go build -a -ldflags '-extldflags "-static"' -o onDefineDomain cmd/main.go

FROM scratch

COPY --from=builder /build/onDefineDomain /usr/bin/onDefineDomain

ENTRYPOINT [ "/usr/bin/onDefineDomain" ]