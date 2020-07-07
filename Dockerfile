FROM golang:alpine

ENV GOPROXY https://mirrors.aliyun.com/goproxy/
ENV PKG_CONFIG_PATH /usr/lib/pkgconfig/ 

RUN mkdir /tmpMake
ADD . /tmpMake
WORKDIR /tmpMake

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update 
RUN apk add --no-cache libzmq-static czmq-dev libsodium-static build-base util-linux-dev
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN go build -o main .

RUN apk add --no-cache tzdata && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
	mkdir -p /server/log && \
	mkdir -p /server/static/fils && \
	cp -r /tmpMake/main /server/ && \
	cp -r /tmpMake/config /server

WORKDIR /server

RUN rm -rf /tmpMake

CMD ["./main"]
