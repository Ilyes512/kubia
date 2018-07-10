FROM golang:1.10.3-alpine3.7 AS builder

WORKDIR $GOPATH/src/app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-w -s' -o kubia . \
    && apk add --upgrade --no-cache \
        ca-certificates

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --from=builder /go/src/app/kubia /kubia

ENTRYPOINT ["/kubia"]
