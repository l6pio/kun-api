FROM golang:1.16 as builder
ENV GOPROXY=https://goproxy.cn
ADD . /src
WORKDIR /src
RUN go mod vendor
RUN go test -v .
RUN gofmt -w .
RUN mkdir -p release
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 CGO_LDFLAGS="-static" go build -mod=vendor -ldflags "-s -w" -v -o release/kun-api .

FROM alpine:3.13 as release
COPY --from=builder /src/release/kun-api /kun-api
COPY bin/grype /usr/local/bin/grype
COPY docker/entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
