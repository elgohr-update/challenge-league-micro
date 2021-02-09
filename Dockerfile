FROM alpine:latest as builder

ENV GOOS linux
ENV GOARCH amd64

RUN apk --no-cache add make git go gcc libtool musl-dev

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

WORKDIR /go/build/micro
COPY . .
COPY --from=localhost:32000/go-micro:dkozlov /go/build/go-micro go-micro
COPY --from=localhost:32000/nakama-go:dkozlov /go/build/nakama-go nakama-go
COPY --from=localhost:32000/nakama-apigrpc:dkozlov /go/build/nakama-apigrpc apigrpc
RUN go mod download


RUN make


FROM alpine:latest 

RUN apk add ca-certificates && \
    rm -rf /var/cache/apk/* /tmp/* && \
    [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf


COPY --from=builder /go/build/micro/micro micro
ENTRYPOINT ["/micro"]
