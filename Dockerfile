FROM golang:1.14.2-alpine3.11 as builder
LABEL maintainer="thien.be@dzones.vn"

ENV GO111MODULE=on \
    CGO_ENABLED=1  \
    GOARCH="amd64" \
    GOOS=linux

RUN apk update && apk add curl make git pkgconfig gcc g++ bash jq \
    gettext
RUN git clone https://github.com/edenhill/librdkafka.git  && \
    cd librdkafka && \
    ./configure --prefix /usr && \
    make && \
    make install
WORKDIR /app
COPY . .
RUN go build -tags musl --ldflags "-extldflags -static" main.go

FROM alpine:3.11
LABEL maintainer="thien.be@dzones.vn"

WORKDIR /home
COPY --from=builder /app/db ./db
COPY --from=builder /app/config/common_bk.config ./config/common.config
COPY --from=builder /app/config/list_ip_bk.config ./config/list_ip.config
COPY --from=builder /app/main .
CMD ["./main"]