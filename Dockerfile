FROM alpine:latest
MAINTAINER Iskakov Zhanat <iskakov_zhanat@mail.ru>
ADD golang-socketio-server /usr/bin/golang-socketio-server
ENTRYPOINT ["golang-socketio-server"]
EXPOSE 5000