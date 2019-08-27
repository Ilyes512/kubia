FROM golang:1.12.9-alpine3.10 AS builder

WORKDIR /src

RUN apk add --no-cache \
        git \
        musl-dev \
        gcc \
        make \
        libc-dev \
        upx

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN mkdir build \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-extldflags "-static" -s -w' -o ./build/kubia

RUN upx --brute --no-progress ./build/kubia

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --from=builder /src/build/kubia /
COPY --from=builder /src/templates /templates
COPY --from=builder /etc/passwd /etc/passwd

CMD ["/kubia"]
