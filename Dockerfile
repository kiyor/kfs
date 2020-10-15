FROM golang:1.15 as builder
WORKDIR /go/src/github.com/kiyor/kfs
COPY go.mod ./
COPY go.sum ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kfs .

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /root
COPY --from=builder /go/src/github.com/kiyor/kfs/kfs .
EXPOSE 8080
ENTRYPOINT ["./kfs"]
