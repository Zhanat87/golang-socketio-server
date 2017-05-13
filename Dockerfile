FROM golang:latest
ADD ~/go/bin/socketio-server /go/socketio-server
ENTRYPOINT /go/socketio-server
EXPOSE 5000