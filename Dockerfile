# build stage
FROM golang:1.15.6 AS builder

WORKDIR /app
COPY . .
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && make clean build

# final stage
FROM alpine
LABEL name=ddns-go url=https://github.com/jeessy2/ddns-go
LABEL maintainer="mingcheng<mingcheng@outlook.com>"

WORKDIR /app
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk add --no-cache tzdata
ENV TZ=Asia/Shanghai
COPY --from=builder /app/ddns-go /app/ddns-go
EXPOSE 9876
ENTRYPOINT /app/ddns-go
