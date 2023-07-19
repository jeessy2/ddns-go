# build stage
FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS builder

WORKDIR /app
COPY . .
ARG TARGETOS TARGETARCH

RUN apk add --no-cache git make tzdata \
    && GOOS=$TARGETOS GOARCH=$TARGETARCH make clean build

# final stage
FROM alpine
LABEL name=ddns-go
LABEL url=https://github.com/jeessy2/ddns-go

WORKDIR /app
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Shanghai
COPY --from=builder /app/ddns-go /app/ddns-go
EXPOSE 9876
ENTRYPOINT ["/app/ddns-go"]
CMD ["-l", ":9876", "-f", "300"]
