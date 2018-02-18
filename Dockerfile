FROM golang:alpine as builder
ARG goos=linux
ARG goarch=amd64
ARG goarm=7
RUN apk update && apk add --update git
RUN mkdir /go/src/app
ADD . /go/src/app/
WORKDIR /go/src/app
ENV CGO_ENABLED=0 GOOS=$goos GOARCH=$goarch GOARM=$goarm
RUN go get -v ./...
RUN cd cmd/podcast-manage-svc && go build -o /podcast-manage-svc

FROM alpine
RUN apk update && apk add --update ca-certificates
COPY --from=builder /podcast-manage-svc /bin/
ENTRYPOINT ["/bin/podcast-manage-svc"]
CMD ["-h"]
