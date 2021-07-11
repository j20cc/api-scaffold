FROM golang:1.16

WORKDIR /go/src/app
COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o bin/api cmd/api/*.go

CMD ["bin/api"]