FROM busybox:latest
COPY ./build/gomsvc /go/bin/gomsvc
ENTRYPOINT [ "/go/bin/gomsvc" ]