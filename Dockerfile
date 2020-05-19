FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
# Create appuser.
ENV USER=appuser
ENV UID=10001 
# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"
WORKDIR $GOPATH/src/github.com/YaleOpenLab/openx
COPY . .
RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o openx
RUN ["cp", "dummyconfig.yaml", "config.yaml"]
RUN ["mv", "config.yaml", "/"]
RUN ["mv", "openx", "/"]
WORKDIR /
RUN ["ls"]
# Step 2: build a smaller image
FROM alpine:3.11
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
WORKDIR /
COPY --from=builder /config.yaml .
COPY --from=builder /openx .
# EXPOSE 8080
RUN ["ls"]
ENTRYPOINT ["/openx"]