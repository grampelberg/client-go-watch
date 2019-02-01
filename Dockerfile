FROM golang:1-alpine AS golang
WORKDIR /go/src/github.com/grampelberg/client-go-watch
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOBIN=/go/bin go install main.go

FROM alpine:3.9
ENV PATH=$PATH:/go/bin
COPY LICENSE /LICENSE
COPY --from=golang /go/bin /go/bin
