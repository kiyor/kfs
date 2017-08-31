FROM golang
ADD . /go/src/github.com/kiyor/kfs
RUN cd /go/src/github.com/kiyor/kfs && \
	go get && \
	go install github.com/kiyor/kfs

EXPOSE 8080
ENTRYPOINT ["/go/bin/kfs"]
