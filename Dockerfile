FROM golang:alpine

RUN mkdir /tmpMake
ADD . /tmpMake
WORKDIR /tmpMake
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