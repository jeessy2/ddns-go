# build stage
FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go get -d -v . \
    && go install -v . \
    && go build -v .

# final stage
FROM alpine
WORKDIR /app
ENV TZ=Asia/Shanghai
COPY --from=builder /app/ddns-go /app/ddns-go
EXPOSE 9876
ENTRYPOINT /app/ddns-go
LABEL Name=ddns-go Version=0.0.1
