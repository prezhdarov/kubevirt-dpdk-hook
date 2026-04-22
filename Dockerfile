FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0

ARG CMD

ADD . /build
WORKDIR /build

RUN go mod download && \
    go build -a -ldflags '-extldflags "-static"' -o app cmd/${CMD}/main.go

FROM scratch

COPY --from=builder /build/app /usr/bin/app

ENTRYPOINT [ "/usr/bin/app" ]
