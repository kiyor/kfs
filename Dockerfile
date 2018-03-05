FROM golang as builder
WORKDIR /go/src/github.com/kiyor/kfs
COPY . .
RUN go get && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kfs .

FROM alpine
WORKDIR /root
COPY --from=builder /go/src/github.com/kiyor/kfs/kfs .
EXPOSE 8080 8081
ENTRYPOINT ["./kfs"]
