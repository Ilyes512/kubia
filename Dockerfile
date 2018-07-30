FROM golang:1.10.3-alpine3.8 AS builder

WORKDIR $GOPATH/src/github.com/Ilyes512/kubia

COPY . .

RUN apk add --no-cache \
        git \
        musl-dev \
        gcc \
        make \
        libc-dev \
        curl \
        xz \
        upx
    
RUN go get -d -v \
    && CC=$(which gcc) GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-linkmode external -extldflags "-static" -s -w' -o kubia

RUN upx --brute --no-progress kubia

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --from=builder /go/src/github.com/Ilyes512/kubia/kubia /
COPY --from=builder /etc/passwd /etc/passwd

ENTRYPOINT ["/kubia"]
