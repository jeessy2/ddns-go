# build stage
FROM golang:1.18 AS builder

WORKDIR /app
COPY . .
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
RUN make clean build

# final stage
FROM scratch
LABEL name=ddns-go
LABEL url=https://github.com/jeessy2/ddns-go

WORKDIR /app
ENV TZ=Asia/Shanghai
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/ddns-go /app/ddns-go
EXPOSE 9876
ENTRYPOINT ["/app/ddns-go"]
CMD ["-l", ":9876", "-f", "300"]
