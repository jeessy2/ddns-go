FROM alpine
LABEL name=ddns-go
LABEL url=https://github.com/jeessy2/ddns-go
RUN apk add --no-cache curl grep

WORKDIR /app
COPY ddns-go /app/
COPY zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Shanghai
EXPOSE 9876
ENTRYPOINT ["/app/ddns-go"]
CMD ["-l", ":9876", "-f", "300"] 