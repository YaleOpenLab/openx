FROM golang
ADD . /go/src/github.com/YaleOpenLab/openx
RUN go get -v github.com/YaleOpenLab/openx
# ENTRYPOINT /go/bin/opensolar
# EXPOSE 8080