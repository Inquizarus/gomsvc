FROM busybox:latest
COPY ./build/gomsvc /go/bin/gomsvc
WORKDIR /app
COPY ./config.json /app/
COPY ./routes /app/routes/
ENV GOMSVC_CONFIG_PATH="/app/config.json"
ENV GOMSVC_ROUTES_DIR="/app/routes"
ENV GOMSVC_LOG_LEVEL="info"
ENTRYPOINT [ "/go/bin/gomsvc" ]