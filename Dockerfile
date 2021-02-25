FROM golang:buster as builder

WORKDIR $GOPATH/src/github.com/davex98/image-clone-controller
COPY . .

RUN go get -d -v
RUN go mod download
RUN go mod verify
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/manager .

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /go/bin/manager .
USER nonroot:nonroot

ENTRYPOINT ["/manager"]
