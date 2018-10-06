# Stage 1. Build the binary
FROM golang:1.11

# add a non-privileged user
RUN useradd -u 10001 myapp

RUN mkdir -p /go/src/github.com/rumyantseva/go-sofia
ADD . /go/src/github.com/rumyantseva/go-sofia
WORKDIR /go/src/github.com/rumyantseva/go-sofia

# build the binary with go build
RUN CGO_ENABLED=0 go build \
	-o bin/go-sofia github.com/rumyantseva/go-sofia/cmd/go-sofia

# Stage 2. Run the binary
FROM scratch

ENV PORT 8080
ENV DIAG_PORT 8585

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=0 /etc/passwd /etc/passwd
USER myapp

COPY --from=0 /go/src/github.com/rumyantseva/go-sofia/bin/go-sofia /go-sofia
EXPOSE $PORT
EXPOSE $DIAG_PORT

CMD ["/go-sofia"]
