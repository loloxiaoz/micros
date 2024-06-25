# Stage 1
FROM docker.af-biz.qianxin-inc.cn/golang:1.20.7 AS builder

ENV GO111MODULE=on \
	GOPROXY=https://goproxy.cn,direct

WORKDIR /go/release

ADD . . 

RUN make build

# Stage 2
FROM docker.af-biz.qianxin-inc.cn/ubuntu:22.04 as prod

RUN apt-get -y update && DEBIAN_FRONTEND="noninteractive" apt -y install tzdata && apt-get clean

ENV LANG='en_US.UTF-8' LANGUAGE='en_US:en' LC_ALL='en_US.UTF-8' TZ=Asia/Shanghai

COPY --from=builder /go/release/build/* ./

EXPOSE 8090/tcp

CMD ["./server", "-env", "fat"]