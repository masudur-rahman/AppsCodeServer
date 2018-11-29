#FROM golang
#
#COPY . /go/src/github/masudur-rahman/AppsCodeServer
#
#RUN go get -u github.com/gorilla/mux
#RUN go get -u github.com/spf13/cobra/cobra
#
#RUN go install /go/src/github/masudur-rahman/AppsCodeServer
#
#ENTRYPOINT ["/go/bin/AppsCodeServer"]
#
#EXPOSE 8080



FROM busybox:glibc

COPY AppsCodeServer /go/bin/AppsCodeServer


EXPOSE 8080

CMD ["--bypass", "true", "--stopTime", "2"]

ENTRYPOINT ["/go/bin/AppsCodeServer", "start"]
